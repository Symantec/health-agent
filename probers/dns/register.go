package dns

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *dnsconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("healthy", &p.healthy, units.None,
		"Is DNS reachable?"); err != nil {
		return err
	}
	latencyBucketer := tricorder.NewGeometricBucketer(0.1, 10e3)
	p.dnsLatencyDistribution = latencyBucketer.NewCumulativeDistribution()
	if err := dir.RegisterMetric("latency", p.dnsLatencyDistribution,
		units.Millisecond, "DNS latency"); err != nil {
		return err
	}
	return nil
}
