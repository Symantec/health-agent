package main

import (
	"errors"
	"os/exec"
	"time"
)

func checkSsh(hostname string, username string, identityFile string) error {
	args := []string{
		"-o", "BatchMode=yes",
		"-o", "CheckHostIP=no",
		"-o", "ConnectionAttempts=1",
		"-o", "ConnectTimeout=5",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	if identityFile != "" {
		args = append(args, "-i", identityFile)
	}
	args = append(args, username+"@"+hostname, "true")
	cmd := exec.Command("ssh", args...)
	err := cmd.Run()
	if err == nil {
		return nil
	}
	return errors.New("Error logging in as: " + username + ": " + err.Error())
}

func checkSshTimeout(hostname, username string, identityFile string,
	stopTime time.Time) error {
	for {
		err := checkSsh(hostname, username, identityFile)
		if err == nil {
			return nil
		}
		if time.Now().Add(*probeInterval).After(stopTime) {
			return err
		}
		time.Sleep(*probeInterval)
	}
}

func runSshCheck(hostname string, username string, identityFile string,
	stopTime time.Time, errorChannel chan<- error) {
	errorChannel <- checkSshTimeout(hostname, username, identityFile, stopTime)
}
