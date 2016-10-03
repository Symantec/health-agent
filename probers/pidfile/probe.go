package pidfile

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func (p *pidconfig) probe() error {
	p.setpidfile(true)
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
	p.process.Pid = int(pid)
	if err := p.process.Signal(syscall.Signal(0)); err != nil {
		p.setpid(false)
		return err
	}
	return nil
}

func (p *pidconfig) sethealthy() {
	p.healthy = p.pidexists && p.pidexists
}

func (p *pidconfig) setpidfile(v bool) {
	p.pidfileexists = v
	p.setpid(v)
}

func (p *pidconfig) setpid(v bool) {
	p.pidexists = v
	p.sethealthy()
}
