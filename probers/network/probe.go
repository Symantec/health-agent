package network

import (
	"os/exec"
	"time"
)

var pingName string = "ping"

func (p *prober) probe() error {
	timeStart := time.Now()
	cmd := exec.Command(pingName, "-c", "1", "-W", "5", p.gatewayAddress)
	if cmd.Run() != nil {
		return nil
	}
	p.gatewayRttDistribution.Add(time.Since(timeStart))
	return nil
}
