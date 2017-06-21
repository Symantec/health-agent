package testprog

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *testprogconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("healthy", &p.healthy,
		units.None, "Healthy?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("exit-code", &p.exitCode,
		units.None, "Program exit value"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("exit-error", &p.exitError,
		units.None, "Probe exited with error"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("stdout", &p.stdout,
		units.None, "Program's stdout"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("stderr", &p.stderr,
		units.None, "Program's stderr"); err != nil {
		return err
	}
	return nil
}
