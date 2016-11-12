package dns

import (
	"net"
	"time"
)

func (p *dnsconfig) probe() error {
	starttime := time.Now()
	if _, err := net.LookupIP(p.hostname); err != nil {
		p.healthy = false
	}
	p.latency = time.Since(starttime) / time.Millisecond
	p.healthy = true
	return nil
}
