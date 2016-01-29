package scheduler

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"time"
)

type loadAverage struct {
	oneMinute      float32
	fiveMinutes    float32
	fifteenMinutes float32
}

type cpuStatistics struct {
	userTime     time.Duration
	userNiceTime time.Duration
	systemTime   time.Duration
	idleTime     time.Duration
	iOWaitTime   time.Duration
	irqTime      time.Duration
	softIrqTime  time.Duration
}

type prober struct {
	loadavg  loadAverage
	numCpus  uint64
	cpuStats cpuStatistics
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
