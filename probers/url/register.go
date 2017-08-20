package url

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *urlconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("has-tricorder-metrics",
		&p.hasTricorderMetrics, units.None, "Tricorder metrics?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("healthy", &p.healthy,
		units.None, "Healthy?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("port-number",
		&p.urlport, units.None, "Port number"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("status-code", &p.statusCode,
		units.None, "Status code"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("error", &p.error,
		units.None, "Error if any"); err != nil {
		return err
	}
	return nil
}
