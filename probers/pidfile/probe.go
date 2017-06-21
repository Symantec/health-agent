package pidfile

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func (p *pidconfig) probe() error {
	if _, err := os.Stat(p.pidfilepath); os.IsNotExist(err) {
		p.setpidfile(false)
		return err
	}
	data, err := ioutil.ReadFile(p.pidfilepath)
	if err != nil {
		p.setpid(false)
		return err
	}
	pidstring := strings.TrimSpace(string(data))
	pid, err := strconv.ParseInt(pidstring, 10, 64)
	if err != nil {
		p.setpid(false)
		return err
	}
	process := os.Process{Pid: int(pid)}
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
