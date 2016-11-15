package dns

import (
	"net"
	"time"
)

func (p *dnsconfig) probe() error {
	starttime := time.Now()
	if _, err := net.LookupIP(p.hostname); err != nil {
		p.healthy = false
	} else {
		p.healthy = true
	}
	p.latency = int64(time.Since(starttime) / time.Millisecond)
	return nil
}
