package main

import (
	"errors"
	"github.com/Symantec/tricorder/go/tricorder/messages"
	"net"
	"net/rpc"
	"time"
)

func checkHealth(address string) ([]string, error) {
	// There is no timeout or context variant for rpc.DialHTTP(), so make a
	// dummy connection to work around this.
	conn, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		return nil, err
	}
	conn.Close()
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

func runHealthCheck(address string, stopTime time.Time,
	errorChannel chan<- error) {
	unhealthyList, err := checkHealthTimeout(address, stopTime)
	if err != nil {
		errorChannel <- errors.New("Error checking health: " + err.Error())
	}
	if len(unhealthyList) > 0 {
		errorChannel <- errors.New(
			address + " has failing health checks: " + err.Error())
	}
	errorChannel <- nil
}
