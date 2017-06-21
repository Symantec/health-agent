package ldap

import (
	"crypto/tls"
	"gopkg.in/ldap.v2"
	"net"
	"time"
)

func (p *ldapconfig) probe() error {
	timeout := time.Duration(p.probefreq) * time.Second
	timeout = timeout / time.Duration(len(p.hostnames))
	p.healthy = false
	for _, hostname := range p.hostnames {
		hostnamePort := hostname + ":636"
		start := time.Now()

		// timeouts must be speficied both at the network layer and at the LDAP layer
		tlsConn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout},
			"tcp", hostnamePort, &tls.Config{ServerName: hostname})
		if err != nil {
			continue
		}
		// we dont close the tls connection directly  close defer to the new ldap connection
		conn := ldap.NewConn(tlsConn, true)
		defer conn.Close()
		conn.SetTimeout(timeout)
		conn.Start()

		err = conn.Bind(p.bindDN, p.bindPassword)
		latency := time.Since(start)
		p.ldapLatencyDistribution.Add(latency)
		if err != nil {
			continue
		}
		p.healthy = true
		break

	}
	return nil
}
