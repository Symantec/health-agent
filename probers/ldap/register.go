package ldap

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *ldapconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("healthy", &p.healthy, units.None,
		"Is LDAP reachable and configured?"); err != nil {
		return err
	}
	latencyBucketer := tricorder.NewGeometricBucketer(0.1, 10e3)
	p.ldapLatencyDistribution = latencyBucketer.NewCumulativeDistribution()
	if err := dir.RegisterMetric("latency", p.ldapLatencyDistribution,
		units.Millisecond, "LDAP latency"); err != nil {
		return err
	}
	return nil
}
