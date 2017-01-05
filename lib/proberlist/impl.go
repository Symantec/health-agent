package proberlist

import (
	"fmt"
	"github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"io"
	"log"
	"time"
)

var (
	latencyBucketer = tricorder.NewGeometricBucketer(0.1, 100e3)
)

func newProberList(proberPath string) *ProberList {
	pl := &ProberList{
		proberPath:            proberPath,
		probeTimeDistribution: latencyBucketer.NewCumulativeDistribution(),
	}
	if err := tricorder.RegisterMetric(proberPath+"/probe-duration",
		pl.probeTimeDistribution, units.Millisecond,
		"duration of last probe"); err != nil {
		panic(err)
	}
	if err := tricorder.RegisterMetric(proberPath+"/probe-start-time",
		&pl.probeStartTime, units.None,
		"start time of last probe"); err != nil {
		panic(err)
	}
	return pl
}

func (pl *ProberList) add(registerProber RegisterProber, path string,
	probeInterval uint8) {
	if registerProber == nil {
		return
	}
	if err := registerProber.Register(mkdir(path)); err != nil {
		panic(err)
	}
	pl.addProber(registerProber, path, probeInterval)
}

func (pl *ProberList) addProber(genericProber prober.Prober, path string,
	probeInterval uint8) {
	newProber := &proberType{
		prober:                genericProber,
		probeInterval:         time.Duration(probeInterval) * time.Second,
		probeTimeDistribution: latencyBucketer.NewCumulativeDistribution(),
	}
	if err := tricorder.RegisterMetric(pl.proberPath+path+"/probe-duration",
		newProber.probeTimeDistribution, units.Millisecond,
		"duration of last probe"); err != nil {
		panic(err)
	}
	if probeInterval > 0 {
		if err := tricorder.RegisterMetric(pl.proberPath+path+"/probe-interval",
			&newProber.probeInterval, units.Second,
			"probe interval"); err != nil {
			panic(err)
		}
		if err := tricorder.RegisterMetric(
			pl.proberPath+path+"/probe-start-time",
			&newProber.probeStartTime, units.None,
			"start time of last probe"); err != nil {
			panic(err)
		}
	}
	pl.lock.Lock()
	defer pl.lock.Unlock()
	pl.probers = append(pl.probers, newProber)
}

func (pl *ProberList) getProbers() []*proberType {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	result := make([]*proberType, len(pl.probers))
	copy(result, pl.probers)
	return result
}

func (pl *ProberList) startProbing(defaultProbeInterval uint,
	logger *log.Logger) {
	for _, p := range pl.getProbers() {
		if p.probeInterval > 0 {
			go p.proberLoop(defaultProbeInterval, logger)
		}
	}
	go pl.proberLoop(defaultProbeInterval, logger)
}

func (pl *ProberList) proberLoop(probeInterval uint, logger *log.Logger) {
	for {
		probeStartTime := time.Now()
		pl.probe(logger)
		probeDuration := time.Since(probeStartTime)
		time.Sleep(time.Second*time.Duration(probeInterval) - probeDuration)
	}
}

func (pl *ProberList) probe(logger *log.Logger) {
	// TODO: Possible data races with tricorder. Assigning a time
	// is not atomic.
	pl.probeStartTime = time.Now()
	for _, p := range pl.getProbers() {
		if p.probeInterval > 0 { // Handled by a dedicated goroutine.
			continue
		}
		startTime := time.Now()
		if err := p.prober.Probe(); err != nil {
			logger.Println(err)
		}
		p.probeTimeDistribution.Add(time.Since(startTime))
	}
	pl.probeTimeDistribution.Add(time.Since(pl.probeStartTime))
}

func (pl *ProberList) writeHtml(writer io.Writer) {
	for _, p := range pl.getProbers() {
		if htmler, ok := p.prober.(HtmlWriter); ok {
			htmler.WriteHtml(writer)
			fmt.Fprintln(writer, "<br>")
		}
	}
}

func (p *proberType) proberLoop(defaultProbeInterval uint, logger *log.Logger) {
	// Set the initial probe interval to the global default, if less than the
	// interval for this prober. The probe interval will be gradually increased
	// until the target probe interval is reached. This gives faster probing at
	// startup when higher resolution may be helpful, and then backs off.
	probeInterval := time.Duration(defaultProbeInterval) * time.Second
	if probeInterval > p.probeInterval {
		probeInterval = p.probeInterval
	}
	for {
		p.probeStartTime = time.Now()
		if err := p.prober.Probe(); err != nil {
			logger.Println(err)
		}
		probeDuration := time.Since(p.probeStartTime)
		p.probeTimeDistribution.Add(probeDuration)
		time.Sleep(probeInterval - probeDuration)
		// Increase the probe interval until the interval for this prober is
		// reached.
		if probeInterval < p.probeInterval {
			probeInterval += time.Duration(defaultProbeInterval) * time.Second
		}
		if probeInterval > p.probeInterval {
			probeInterval = p.probeInterval
		}
	}
}

func mkdir(name string) *tricorder.DirectorySpec {
	dir, err := tricorder.RegisterDirectory(name)
	if err != nil {
		panic(err)
	}
	return dir
}
