package testprog

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type testprogconfig struct {
	testname  string
	filepath  string
	healthy   bool
	exitCode  int
	exitError string
	stdout    string
	stderr    string
}

func (p *testprogconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *testprogconfig) Probe() error {
	return p.probe()
}

func Maketestprogprober(testname string, testprogpath string) *testprogconfig {
	p := new(testprogconfig)
	p.testname = testname
	p.filepath = testprogpath
	return p
}
