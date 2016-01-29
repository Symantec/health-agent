package scheduler

import (
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
)

// TODO(rgooch): Figure out how to share with systime package.
var onlineCpuFilename string = "/sys/devices/system/cpu/online"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	p.numCpus = getNumCpus()
	p.registerLoadavg(dir)
	p.registerCpu(dir)
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

func mkdir(dir *tricorder.DirectorySpec, name string) *tricorder.DirectorySpec {
	subdir, err := dir.RegisterDirectory(name)
	if err != nil {
		panic(err)
	}
	return subdir
}

func (p *prober) registerLoadavg(dir *tricorder.DirectorySpec) error {
	dir = mkdir(dir, "loadavg")
	if err := dir.RegisterMetric("1m", &p.loadavg.oneMinute, units.None,
		"load average for the last minute"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("5m", &p.loadavg.fiveMinutes,
		units.None, "load average for the last minute"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("15m", &p.loadavg.fifteenMinutes,
		units.None, "load average for the last minute"); err != nil {
		return err
	}
	return nil
}

func (p *prober) registerCpu(dir *tricorder.DirectorySpec) error {
	dir = mkdir(dir, "cpu")
	if err := p.cpuStats.idleTime.register(dir, "idle-time",
		"idle time"); err != nil {
		return err
	}
	if err := p.cpuStats.iOWaitTime.register(dir, "iowait-time",
		"time spent waiting for I/O to complete"); err != nil {
		return err
	}
	if err := p.cpuStats.irqTime.register(dir, "irq-time",
		"time spent servicing interrupts"); err != nil {
		return err
	}
	if err := p.cpuStats.userNiceTime.register(dir, "nice-time",
		"niced processes executing in user mode"); err != nil {
		return err
	}
	if err := p.cpuStats.softIrqTime.register(dir, "softirq-time",
		"time spent servicing softirqs"); err != nil {
		return err
	}
	if err := p.cpuStats.systemTime.register(dir, "system-time",
		"processes executing in kernel mode"); err != nil {
		return err
	}
	if err := p.cpuStats.userTime.register(dir, "user-time",
		"normal priority processes executing in user mode"); err != nil {
		return err
	}
	return nil
}

func (m *timeMetric) register(dir *tricorder.DirectorySpec,
	name, use string) error {
	if err := dir.RegisterMetric(name, &m.value, units.None, use); err != nil {
		return err
	}
	return dir.RegisterMetric(name+"-fraction", &m.fraction, units.None,
		use+" during last probe interval")
}
