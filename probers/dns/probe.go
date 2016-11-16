package dns

import (
	"net"
	"time"
)

func (p *dnsconfig) probe() error {
	for hostname := range p.hostnames {
		start := time.Now()
		if _, err := net.LookupIP(p.hostname); err != nil {
			p.healthy = false
			p.latency = time.Since(start)
			continue
		} else {
			p.healthy = true
			p.latency = time.Since(start)
			return nil
		}
	}
	return nil
}
