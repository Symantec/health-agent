package main

import (
	"flag"
	"fmt"
	"github.com/Symantec/Dominator/lib/logbuf"
	"github.com/Symantec/health-agent/httpd"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"log"
	"net/rpc"
	"os"
	"runtime"
	"time"
)

var (
	logbufLines = flag.Uint("logbufLines", 1024,
		"Number of lines to store in the log buffer")
	maxThreads = flag.Uint("maxThreads", 1,
		"Maximum number of parallel OS threads to use")
	portNum = flag.Uint("portNum", 6910,
		"Port number to allocate and listen on for HTTP/RPC")
	probeInterval = flag.Uint("probeInterval", 10, "Probe interval in seconds")
)

func doMain() error {
	flag.Parse()
	runtime.GOMAXPROCS(int(*maxThreads))
	runtime.LockOSThread()
	circularBuffer := logbuf.New(*logbufLines)
	logger := log.New(circularBuffer, "", log.LstdFlags)
	probers, err := setupProbers()
	if err != nil {
		return err
	}
	latencyBucketer := tricorder.NewGeometricBucketer(0.1, 100e3)
	scanTimeDistribution := latencyBucketer.NewDistribution()
	var scanStartTime time.Time
	if err := tricorder.RegisterMetric("scan-duration",
		scanTimeDistribution, units.Millisecond,
		"duration of last probe"); err != nil {
		return err
	}
	if err := tricorder.RegisterMetric("scan-start-time", &scanStartTime,
		units.None, "start time of last probe"); err != nil {
		return err
	}
	rpc.HandleHTTP()
	httpd.AddHtmlWriter(circularBuffer)
	if err := httpd.StartServer(*portNum); err != nil {
		return err
	}
	for {
		scanStartTime = time.Now()
		for _, p := range probers {
			if err := p.Probe(); err != nil {
				logger.Println(err)
			}
		}
		scanDuration := time.Since(scanStartTime)
		scanTimeDistribution.Add(scanDuration)
		time.Sleep(time.Second*time.Duration(*probeInterval) - scanDuration)
	}
	_ = logger
	return nil
}

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
