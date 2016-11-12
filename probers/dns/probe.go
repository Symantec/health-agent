package dns

import (
	"net"
	"time"
)

func (p *dnsconfig) probe() error {
	hostname := p.hostname
	starttime := time.Now()
	conn, err := net.Dial("tcp", hostname)
	p.latency = time.Since(starttime).String()
	if err != nil {
		p.healthy = false
		p.err = err
	}
	defer conn.Close()
	p.healthy = true
	p.err = nil
	return nil
}
