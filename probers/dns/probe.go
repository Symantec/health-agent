package dns

import (
	"net"
	"time"
)

func (p *dnsconfig) probe() error {
	start := time.Now()
	if _, err := net.LookupIP(p.hostname); err != nil {
		p.healthy = false
	} else {
		p.healthy = true
	}
	latency := time.Since(start)
	p.dnsLatencyDistribution.Add(latency)
	return nil
}
