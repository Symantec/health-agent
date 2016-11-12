package dns

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type dnsconfig struct {
	testname string
	hostname string
	healthy  bool
	latency  int64
}

func (p *dnsconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *dnsconfig) Probe() error {
	return p.probe()
}

func Makednsprober(testname, hostname string) *dnsconfig {
	p := new(dnsconfig)
	p.testname = testname
	p.hostname = hostname
	return p
}
