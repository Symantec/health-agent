package pidfile

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *pidconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("healthy", &p.healthy, units.None,
		"Is process healthy?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("pidfile-exists", &p.pidFileExists, units.None,
		"Does pidfile exist?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("pid-exists", &p.pidExists, units.None,
		"Does pid exist?"); err != nil {
		return err
	}
	return nil
}
