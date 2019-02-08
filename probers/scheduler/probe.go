package scheduler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	loadavgFilename = "/proc/loadavg"
	statFilename    = "/proc/stat"
)

func (p *prober) probe() error {
	if err := p.probeLoadavg(); err != nil {
		return err
	}
	if err := p.probeCpu(); err != nil {
		return err
	}
	return nil
}

func (p *prober) probeLoadavg() error {
	file, err := os.Open(loadavgFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	var str string
	nScanned, err := fmt.Fscanf(file, "%f %f %f %s",
		&p.loadavg.oneMinute, &p.loadavg.fiveMinutes, &p.loadavg.fifteenMinutes,
		&str)
	if err != nil {
		return err
	}
	if nScanned < 3 {
		return fmt.Errorf("only read %d values from %s",
			nScanned, loadavgFilename)
	}
	return nil
}

func (p *prober) probeCpu() error {
	file, err := os.Open(statFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentTime := time.Now()
	lastCpuStats := p.cpuStats
	for scanner.Scan() {
		if err := p.processStatLine(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !p.cpuStats.lastProbeTime.IsZero() {
		interval := float64(currentTime.Sub(p.cpuStats.lastProbeTime))
		p.cpuStats.idleTime.fraction = float64(p.cpuStats.idleTime.value-
			lastCpuStats.idleTime.value) / interval
		p.cpuStats.iOWaitTime.fraction = float64(p.cpuStats.iOWaitTime.value-
			lastCpuStats.iOWaitTime.value) / interval
		p.cpuStats.irqTime.fraction = float64(p.cpuStats.irqTime.value-
			lastCpuStats.irqTime.value) / interval
		p.cpuStats.userNiceTime.fraction = float64(
			p.cpuStats.userNiceTime.value-lastCpuStats.userNiceTime.value) /
			interval
		p.cpuStats.softIrqTime.fraction = float64(p.cpuStats.softIrqTime.value-
			lastCpuStats.softIrqTime.value) / interval
		p.cpuStats.systemTime.fraction = float64(p.cpuStats.systemTime.value-
			lastCpuStats.systemTime.value) / interval
		p.cpuStats.userTime.fraction = float64(p.cpuStats.userTime.value-
			lastCpuStats.userTime.value) / interval
	}
	p.cpuStats.lastProbeTime = currentTime
	return nil
}

func (p *prober) processStatLine(line string) error {
	if strings.HasPrefix(line, "cpu ") {
		var user, niced, sys, idle, iowait, irq, softIrq int64
		nScanned, err := fmt.Sscanf(line[4:], "%d %d %d %d %d %d %d",
			&user, &niced, &sys, &idle, &iowait, &irq, &softIrq)
		if err != nil {
			return err
		}
		if nScanned < 7 {
			return fmt.Errorf("only read %d values from %s",
				nScanned, statFilename)
		}
		p.cpuStats.idleTime.value = p.tickToDuration(idle)
		p.cpuStats.iOWaitTime.value = p.tickToDuration(iowait)
		p.cpuStats.irqTime.value = p.tickToDuration(irq)
		p.cpuStats.userNiceTime.value = p.tickToDuration(niced)
		p.cpuStats.softIrqTime.value = p.tickToDuration(softIrq)
		p.cpuStats.systemTime.value = p.tickToDuration(sys)
		p.cpuStats.userTime.value = p.tickToDuration(user)
	}
	return nil
}

func (p *prober) tickToDuration(tick int64) time.Duration {
	// TODO(rgooch): Use sysconf(_SC_CLK_TCK).
	tickToDuration := time.Millisecond * 10 / time.Duration(p.numCpus)
	return time.Duration(tick) * tickToDuration
}
