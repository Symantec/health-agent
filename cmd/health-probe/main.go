package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder/messages"
	"net/rpc"
	"os"
	"time"
)

var (
	hostname = flag.String("hostname", "localhost",
		"Hostname of machine to probe")
	portNum       = flag.Uint("portNum", 6910, "Port number of health-agent")
	probeInterval = flag.Duration("probeInterval", time.Second*5,
		"Time between probe intervals (min 100 milliseconds)")
	timeout = flag.Duration("timeout", time.Minute*5, "Time before giving up")
)

func printUsage() {
	fmt.Fprintln(os.Stderr,
		"Usage: health-probe [flags...]")
	fmt.Fprintln(os.Stderr, "Common flags:")
	flag.PrintDefaults()
}

func checkHealth(address string) ([]string, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	var metric messages.Metric
	err = client.Call("MetricsServer.GetMetric", "/health-checks/*/healthy",
		&metric)
	if err != nil {
		return nil, err
	}
	if healthy, ok := metric.Value.(bool); !ok {
		return nil, errors.New("metric value is not bool")
	} else if healthy {
		return nil, nil
	}
	err = client.Call("MetricsServer.GetMetric",
		"/health-checks/*/unhealthy-list", &metric)
	if err != nil {
		return nil, err
	}
	if list, ok := metric.Value.([]string); !ok {
		return nil, errors.New("list metric is not []string")
	} else {
		return list, nil
	}
}

func checkHealthTimeout(address string, stopTime time.Time) ([]string, error) {
	for {
		unhealthyList, err := checkHealth(address)
		if len(unhealthyList) < 1 && err == nil {
			return nil, nil
		}
		if time.Now().Add(*probeInterval).After(stopTime) {
			return unhealthyList, err
		}
		time.Sleep(*probeInterval)
	}
}

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *hostname == "" {
		fmt.Fprintf(os.Stderr, "No hostname specified\n")
		os.Exit(2)
	}
	if *probeInterval < time.Millisecond*100 {
		fmt.Fprintf(os.Stderr, "probeInterval too short\n")
		os.Exit(2)
	}
	address := fmt.Sprintf("%s:%d", *hostname, *portNum)
	unhealthyList, err := checkHealthTimeout(address, time.Now().Add(*timeout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking health: %s\n", err)
		os.Exit(1)
	}
	if len(unhealthyList) > 0 {
		fmt.Fprintf(os.Stderr, "%s has failing health checks: %s\n",
			*hostname, unhealthyList)
		os.Exit(1)
	}
}
