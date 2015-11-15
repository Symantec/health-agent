package systime

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"time"
)

type prober struct {
	dir       *tricorder.DirectorySpec
	idleTime  time.Duration
	probeTime time.Time
	upTime    time.Duration
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
