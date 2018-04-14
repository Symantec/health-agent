package network

import (
	libprober "github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	gatewayAddress              string
	gatewayInterfaceName        string
	gatewayPingTimeDistribution *tricorder.CumulativeDistribution
	gatewayRttDistribution      *tricorder.CumulativeDistribution
	resolverDomain              string
	resolverNameservers         *tricorder.List
	resolverSearchDomains       *tricorder.List
}

func Register(dir *tricorder.DirectorySpec) libprober.Prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
