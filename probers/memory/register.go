package memory

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := dir.RegisterMetric("available", &p.available, units.Byte,
		"estimate of memory available for applications"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("free", &p.free, units.Byte,
		"free memory"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("total", &p.total, units.Byte,
		"total memory"); err != nil {
		panic(err)
	}
	return p
}
