package loadavg

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir            *tricorder.DirectorySpec
	oneMinute      float32
	fiveMinutes    float32
	fifteenMinutes float32
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
