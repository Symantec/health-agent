package ldap

import (
	"github.com/nmcclain/ldap"
	"net"
	"time"
)

func (p *ldapconfig) probe() error {
	timeout := time.Duration(p.probefreq) * time.Second
	timeout = timeout / time.Duration(len(p.hostnames))
	for _, hostname := range p.hostnames {
		hostnamePort := hostname + ":636"
		start := time.Now()
		conn, err := ldap.DialTLSDialer("tcp", hostnamePort,
			nil, &net.Dialer{Timeout: timeout})
		if err != nil {
			p.healthy = false
			continue
		}
		defer conn.Close()
		err = conn.Bind(p.bindDN, p.bindPassword)
		latency := time.Since(start)
		p.ldapLatencyDistribution.Add(latency)
		if err != nil {
			p.healthy = false
		} else {
			p.healthy = true
			break
		}
	}
	return nil
}
