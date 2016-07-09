package kernel

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
