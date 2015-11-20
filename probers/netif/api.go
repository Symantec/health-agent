package netif

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"io"
	"time"
)

type prober struct {
	dir           *tricorder.DirectorySpec
	netInterfaces map[string]*netInterface // map[name]*netInterface
}

type netInterface struct {
	dir                 *tricorder.DirectorySpec
	name                string
	carrier             bool
	multicastFrames     uint64
	rxCompressedPackets uint64
	rxData              uint64
	rxDropped           uint64
	rxErrors            uint64
	rxFrameErrors       uint64
	rxOverruns          uint64
	rxPackets           uint64
	speed               uint64
	txCarrierLosses     uint64
	txCollisionErrors   uint64
	txCompressedPackets uint64
	txData              uint64
	txDropped           uint64
	txErrors            uint64
	txOverruns          uint64
	txPackets           uint64
	probed              bool
	lastProbeTime       time.Time
	lastRxData          uint64
	rxDataRate          uint64
	lastTxData          uint64
	txDataRate          uint64
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}

func (p *prober) WriteHtml(writer io.Writer) {
	p.writeHtml(writer)
}
