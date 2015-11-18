package filesystems

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir         *tricorder.DirectorySpec
	fileSystems map[string]*fileSystem // map[device]*fileSystem
}

type fileSystem struct {
	dir       *tricorder.DirectorySpec
	available uint64
	device    string
	free      uint64
	options   string
	size      uint64
	writable  bool
	probed    bool
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
