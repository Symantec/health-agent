package scheduler

import (
	"bufio"
	"errors"
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
		return errors.New(fmt.Sprintf("only read %d values from %s",
			nScanned, loadavgFilename))
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
	for scanner.Scan() {
		if err := p.processStatLine(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
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
		p.cpuStats.idleTime = p.tickToDuration(idle)
		p.cpuStats.iOWaitTime = p.tickToDuration(iowait)
		p.cpuStats.irqTime = p.tickToDuration(irq)
		p.cpuStats.userNiceTime = p.tickToDuration(niced)
		p.cpuStats.softIrqTime = p.tickToDuration(softIrq)
		p.cpuStats.systemTime = p.tickToDuration(sys)
		p.cpuStats.userTime = p.tickToDuration(user)
	}
	return nil
}

func (p *prober) tickToDuration(tick int64) time.Duration {
	// TODO(rgooch): Use sysconf(_SC_CLK_TCK).
	tickToDuration := time.Millisecond * 10 / time.Duration(p.numCpus)
	return time.Duration(tick) * tickToDuration
}
