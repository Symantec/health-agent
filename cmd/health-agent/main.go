package main

import (
	"flag"
	"fmt"
	"github.com/Symantec/Dominator/lib/logbuf"
	"github.com/Symantec/health-agent/httpd"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	configDir = flag.String("configDir", "/etc/health-agent",
		"Name of the directory containing the health check configs")
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
		return
	}
	defer file.Close()
	fmt.Fprintln(file, os.Getpid())
}

func doMain() error {
	flag.Parse()
	runtime.GOMAXPROCS(int(*maxThreads))
	runtime.LockOSThread()
	circularBuffer := logbuf.New()
	logger := log.New(circularBuffer, "", log.LstdFlags)
	proberList, err := setupProbers()
	if err != nil {
		return err
	}
	if err := setupHealthchecks(*configDir, proberList, logger); err != nil {
		logger.Printf("Error occured while setting up Healthchecks")
		return err
	}
	httpd.AddHtmlWriter(proberList)
	httpd.AddHtmlWriter(circularBuffer)
	sighupChannel := make(chan os.Signal)
	signal.Notify(sighupChannel, syscall.SIGHUP)
	sigtermChannel := make(chan os.Signal)
	signal.Notify(sigtermChannel, syscall.SIGTERM, syscall.SIGINT)
	rpc.HandleHTTP()
	if err := httpd.StartServer(*portNum); err != nil {
		return err
	}
	writePidfile()
	proberList.StartProbing(*probeInterval, logger)
	for {
		select {
		case <-sighupChannel:
			err = syscall.Exec(os.Args[0], os.Args, os.Environ())
			if err != nil {
				logger.Printf("Unable to Exec:%s\t%s\n", os.Args[0], err)
			}
		case <-sigtermChannel:
			gracefulCleanup()
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
