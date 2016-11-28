package ldap

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type ldapconfig struct {
	testname                string
	probefreq               uint8
	hostnames               []string
	bindDN                  string
	bindPassword            string
	healthy                 bool
	ldapLatencyDistribution *tricorder.CumulativeDistribution
}

func (p *ldapconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *ldapconfig) Probe() error {
	return p.probe()
}

func New(testname string, probefreq uint8, hostnames []string,
	binddn string, bindpwd string) *ldapconfig {
	return &ldapconfig{testname: testname,
		probefreq:    probefreq,
		hostnames:    hostnames,
		bindDN:       binddn,
		bindPassword: bindpwd}
}
