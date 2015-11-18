package loadavg

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := dir.RegisterMetric("1m", &p.oneMinute, units.None,
		"load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("5m", &p.fiveMinutes, units.None,
		"load average for the last minute"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("15m", &p.fifteenMinutes, units.None,
		"load average for the last minute"); err != nil {
		panic(err)
	}
	return p
}
