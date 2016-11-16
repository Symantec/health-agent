package dns

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"time"
)

type dnsconfig struct {
	testname string
	hostname string
	healthy  bool
	latency  time.Duration
}

func (p *dnsconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *dnsconfig) Probe() error {
	return p.probe()
}

func New(testname, hostname string) *dnsconfig {
	return &dnsconfig{testname: testname, hostname: hostname}

}
