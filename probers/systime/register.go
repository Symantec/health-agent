package systime

import (
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
)

var onlineCpuFilename string = "/sys/devices/system/cpu/online"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	p.numCpus = getNumCpus()
	if p.numCpus > 0 {
		if err := dir.RegisterMetric("idle-time", &p.idleTime, units.Second,
			"idle time since last boot"); err != nil {
			panic(err)
		}
	}
	getNumCpus()
	if err := dir.RegisterMetric("time", &p.probeTime, units.None,
		"time of last probe"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("uptime", &p.upTime, units.Second,
		"time since last boot"); err != nil {
		panic(err)
	}
	return p
}

func getNumCpus() uint64 {
	file, err := os.Open(onlineCpuFilename)
	if err != nil {
		return 0
	}
	defer file.Close()
	var low, high uint64
	nScanned, err := fmt.Fscanf(file, "%d-%d", &low, &high)
	if err != nil || nScanned != 2 {
		return 0
	}
	return high - low + 1
}
