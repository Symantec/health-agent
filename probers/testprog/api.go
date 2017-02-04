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

func Maketestprogprober(testname string, testprogpath string) *testprogconfig {
	p := new(testprogconfig)
	p.testname = testname
	p.filepath = testprogpath
	return p
}

func (p *testprogconfig) HealthCheck() bool {
	return p.healthy
}

func (p *testprogconfig) Probe() error {
	return p.probe()
}

func (p *testprogconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}
