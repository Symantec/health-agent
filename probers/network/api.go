package network

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	gatewayAddress              string
	gatewayInterfaceName        string
	gatewayPingTimeDistribution *tricorder.CumulativeDistribution
	gatewayRttDistribution      *tricorder.CumulativeDistribution
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
