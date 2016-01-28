package scheduler

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type loadAverage struct {
	oneMinute      float32
	fiveMinutes    float32
	fifteenMinutes float32
}

type prober struct {
	loadavg loadAverage
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
