package memory

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	available uint64
	free      uint64
	total     uint64
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
