package netif

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir           *tricorder.DirectorySpec
	netInterfaces map[string]*netInterface // map[name]*netInterface
}

type netInterface struct {
	dir       *tricorder.DirectorySpec
	rxData    uint64
	rxPackets uint64
	txData    uint64
	txPackets uint64
	probed    bool
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
