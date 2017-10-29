package pidfile

import (
	"fmt"
	"os"
	"syscall"
)

func (p *pidconfig) probe() error {
	if _, err := os.Stat(p.pidfilepath); os.IsNotExist(err) {
		p.setpidfile(false)
		return err
	}
	file, err := os.Open(p.pidfilepath)
	if err != nil {
		p.setpid(false)
		return err
	}
	defer file.Close()
	var pid int
	if nScanned, err := fmt.Fscanf(file, "%d", &pid); err != nil {
		p.setpid(false)
		return err
	} else if nScanned != 1 {
		p.setpid(false)
		return err
	}
	process := os.Process{Pid: pid}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		p.setpid(false)
		return err
	}
	p.setpidfile(true)
	return nil
}

func (p *pidconfig) sethealthy() {
	p.healthy = p.pidFileExists && p.pidExists
}

func (p *pidconfig) setpidfile(v bool) {
	p.pidFileExists = v
	p.setpid(v)
}

func (p *pidconfig) setpid(v bool) {
	p.pidExists = v
	p.sethealthy()
}
