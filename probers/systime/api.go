package systime

import (
	libprober "github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
	"time"
)

type prober struct {
	numCpus   uint64
	idleTime  time.Duration
	probeTime time.Time
	upTime    time.Duration
}

func Register(dir *tricorder.DirectorySpec) libprober.Prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
