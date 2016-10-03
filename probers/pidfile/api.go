package pidfile

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"os"
)

type pidconfig struct {
	processname   string
	pidfilepath   string
	process       os.Process
	healthy       bool
	pidfileexists bool
	pidexists     bool
}

func (p *pidconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *pidconfig) Probe() error {
	return p.probe()
}

func Makepidprober(testname string, pidpath string) *pidconfig {
	p := new(pidconfig)
	p.processname = testname
	p.pidfilepath = pidpath
	return p
}
