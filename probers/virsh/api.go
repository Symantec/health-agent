package virsh

import (
	"github.com/Symantec/health-agent/lib/proberlist"
	"github.com/Symantec/tricorder/go/tricorder"
)

type domainInfo struct {
	dir   *tricorder.DirectorySpec
	state string
}

type prober struct {
	domainsDir  *tricorder.DirectorySpec
	listResults map[string]string // Key: domain name, value: state.
	domains     map[string]*domainInfo
}

func New() proberlist.RegisterProber {
	return newProber()
}

func (p *prober) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
