package script

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type scriptconfig struct {
	scriptname	string
	scriptfilepath	string
	runSuccessful	bool
	exitError	string
}

func (p *scriptconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *scriptconfig) Probe() error {
	return p.probe()
}

func Makescriptprober(testname string, scriptpath string) *scriptconfig {
	p := new(scriptconfig)
	p.scriptname = testname
	p.scriptfilepath = scriptpath
	return p
}
