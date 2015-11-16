package network

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	gatewayInterface       string
	gatewayRttDistribution *tricorder.Distribution
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
