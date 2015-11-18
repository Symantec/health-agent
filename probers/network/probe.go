package network

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var pingName string = "ping"

func (p *prober) probe() error {
	cmd := exec.Command(pingName, "-c", "1", "-W", "5", p.gatewayAddress)
	timeStart := time.Now()
	stdout, err := cmd.Output()
	if err != nil {
		return nil
	}
	pingTime := time.Since(timeStart)
	p.gatewayPingTimeDistribution.Add(pingTime)
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")
		for index, field := range splitLine {
			if strings.HasPrefix(field, "time=") &&
				index < len(splitLine)-1 &&
				splitLine[index+1] == "ms" {
				var rtt float64
				nScanned, err := fmt.Sscanf(field[5:], "%f", &rtt)
				if nScanned == 1 && err == nil {
					p.gatewayRttDistribution.Add(rtt)
					return nil
				}
			}
		}
	}
	// Unable to parse output from ping(8), so use the ping time.
	p.gatewayRttDistribution.Add(pingTime)
	return nil
}
