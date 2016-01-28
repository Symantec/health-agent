package scheduler

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := dir.RegisterMetric("loadavg/1m", &p.loadavg.oneMinute, units.None,
		"load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("loadavg/5m", &p.loadavg.fiveMinutes,
		units.None, "load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("loadavg/15m", &p.loadavg.fifteenMinutes,
		units.None, "load average for the last minute"); err != nil {
		panic(err)
	}
	return p
}
