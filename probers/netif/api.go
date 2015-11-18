package netif

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir           *tricorder.DirectorySpec
	netInterfaces map[string]*netInterface // map[name]*netInterface
}

type netInterface struct {
	dir                 *tricorder.DirectorySpec
	multicastFrames     uint64
	rxCompressedPackets uint64
	rxData              uint64
	rxDropped           uint64
	rxErrors            uint64
	rxFrameErrors       uint64
	rxOverruns          uint64
	rxPackets           uint64
	txCarrierLosses     uint64
	txCollisionErrors   uint64
	txCompressedPackets uint64
	txData              uint64
	txDropped           uint64
	txErrors            uint64
	txOverruns          uint64
	txPackets           uint64
	probed              bool
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
