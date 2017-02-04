package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	hostname = flag.String("hostname", "localhost",
		"Hostname of machine to probe")
	portNum       = flag.Uint("portNum", 6910, "Port number of health-agent")
	probeInterval = flag.Duration("probeInterval", time.Second*5,
		"Time between probe intervals (min 100 milliseconds)")
	sshIdentityFile = flag.String("sshIdentityFile", "",
		"Optional SSH identify file to use")
	timeout = flag.Duration("timeout", time.Minute*5, "Time before giving up")
)

func printUsage() {
	fmt.Fprintln(os.Stderr,
		"Usage: health-probe [flags...] ssh-user...")
	fmt.Fprintln(os.Stderr, "Common flags:")
	flag.PrintDefaults()
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
	stopTime := time.Now().Add(*timeout)
	errorChannel := make(chan error, 0)
	go runHealthCheck(address, stopTime, errorChannel)
	numToHarvest := 1
	for _, username := range flag.Args() {
		go runSshCheck(*hostname, username, *sshIdentityFile, stopTime,
			errorChannel)
		numToHarvest++
	}
	checkFailed := false
	for i := 0; i < numToHarvest; i++ {
		if err := <-errorChannel; err != nil {
			fmt.Fprintln(os.Stderr, err)
			checkFailed = true
		}
	}
	if checkFailed {
		os.Exit(1)
	}
}
