package systime

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	p.dir = dir
	// TODO(rgooch): Consider dividing this by the number of CPUs before
	//               exporting.
	//if err := dir.RegisterMetric("idle-time", &p.idleTime, units.Second,
	//	"idle time since last boot"); err != nil {
	//	panic(err)
	//}
	if err := dir.RegisterMetric("time", &p.probeTime, units.Second,
		"time of last probe"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("uptime", &p.upTime, units.Second,
		"time since last boot"); err != nil {
		panic(err)
	}
	return p
}
