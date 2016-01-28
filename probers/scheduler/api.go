package scheduler

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	loadavgOneMinute      float32
	loadavgFiveMinutes    float32
	loadavgFifteenMinutes float32
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
