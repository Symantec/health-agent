package main

import (
	"errors"
	"os/exec"
	"time"
)

func checkSsh(hostname string, username string) error {
	cmd := exec.Command("ssh",
		"-o", "BatchMode=yes",
		"-o", "CheckHostIP=no",
		"-o", "ConnectionAttempts=1",
		"-o", "ConnectTimeout=5",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		username+"@"+hostname, "true")
	err := cmd.Run()
	if err == nil {
		return nil
	}
	return errors.New("Error logging in as: " + username + ": " + err.Error())
}

func checkSshTimeout(hostname, username string, stopTime time.Time) error {
	for {
		err := checkSsh(hostname, username)
		if err == nil {
			return nil
		}
		if time.Now().Add(*probeInterval).After(stopTime) {
			return err
		}
		time.Sleep(*probeInterval)
	}
}

func runSshCheck(hostname string, username string, stopTime time.Time,
	errorChannel chan<- error) {
	errorChannel <- checkSshTimeout(hostname, username, stopTime)
}
