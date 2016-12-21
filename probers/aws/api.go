package aws

import (
	libprober "github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	// left empty for now
}

func Register(dir *tricorder.DirectorySpec) libprober.Prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
