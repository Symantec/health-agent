package testprog

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type scriptconfig struct {
	testname  string
	filepath  string
	healthy   bool
	exitCode  int
	exitError string
	stdout    string
	stderr    string
}

func (p *scriptconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *scriptconfig) Probe() error {
	return p.probe()
}

func Makescriptprober(testname string, scriptpath string) *scriptconfig {
	p := new(scriptconfig)
	p.testname = testname
	p.filepath = scriptpath
	return p
}
