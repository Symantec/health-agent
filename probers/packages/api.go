package packages

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type packageEntry struct {
	dir     *tricorder.DirectorySpec
	name    string
	version string
	size    uint64
}

type packageList struct {
	dir      *tricorder.DirectorySpec
	packages map[string]*packageEntry
}

type prober struct {
	dir    *tricorder.DirectorySpec
	debian *packageList
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
