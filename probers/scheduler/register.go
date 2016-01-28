package scheduler

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := dir.RegisterMetric("loadavg/1m", &p.loadavgOneMinute, units.None,
		"load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("loadavg/5m", &p.loadavgFiveMinutes,
		units.None, "load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("loadavg/15m", &p.loadavgFifteenMinutes,
		units.None, "load average for the last minute"); err != nil {
		panic(err)
	}
	return p
}
