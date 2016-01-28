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
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	logbufLines = flag.Uint("logbufLines", 1024,
		"Number of lines to store in the log buffer")
	maxThreads = flag.Uint("maxThreads", 1,
		"Maximum number of parallel OS threads to use")
	pidfile = flag.String("pidfile", "/var/run/health-agent.pid",
		"Name of file to write my PID to")
	portNum = flag.Uint("portNum", 6910,
		"Port number to allocate and listen on for HTTP/RPC")
	probeInterval = flag.Uint("probeInterval", 10, "Probe interval in seconds")
)

func gracefulCleanup() {
	if *pidfile != "" {
		os.Remove(*pidfile)
	}
	os.Exit(1)
}

func writePidfile() {
	if *pidfile == "" {
		return
	}
	file, err := os.Create(*pidfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	fmt.Fprintln(file, os.Getpid())
}

func doMain() error {
	flag.Parse()
	runtime.GOMAXPROCS(int(*maxThreads))
	runtime.LockOSThread()
	circularBuffer := logbuf.New(*logbufLines)
	logger := log.New(circularBuffer, "", log.LstdFlags)
	proberList, err := setupProbers()
	if err != nil {
		return err
	}
	latencyBucketer := tricorder.NewGeometricBucketer(0.1, 100e3)
	scanTimeDistribution := latencyBucketer.NewCumulativeDistribution()
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
	httpd.AddHtmlWriter(proberList)
	httpd.AddHtmlWriter(circularBuffer)
	if err := httpd.StartServer(*portNum); err != nil {
		return err
	}
	sighupChannel := make(chan os.Signal)
	signal.Notify(sighupChannel, syscall.SIGHUP)
	sigtermChannel := make(chan os.Signal)
	signal.Notify(sigtermChannel, syscall.SIGTERM, syscall.SIGINT)
	startProbesChannel := make(chan bool, 1)
	writePidfile()
	startProbesChannel <- true
	for {
		select {
		case <-sighupChannel:
			err = syscall.Exec(os.Args[0], os.Args, os.Environ())
			if err != nil {
				logger.Printf("Unable to Exec:%s\t%s\n", os.Args[0], err)
			}
		case <-sigtermChannel:
			gracefulCleanup()
		case <-startProbesChannel:
			scanStartTime = time.Now()
			proberList.Probe(logger)
			scanDuration := time.Since(scanStartTime)
			scanTimeDistribution.Add(scanDuration)
			go func(sleepDuration time.Duration) {
				time.Sleep(sleepDuration)
				startProbesChannel <- true
			}(time.Second*time.Duration(*probeInterval) - scanDuration)
		}
	}
	return nil
}

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
