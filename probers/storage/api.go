package storage

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir            *tricorder.DirectorySpec
	storageDevices map[string]*storageDevice // map[name]*storageDevice
}

type storageDevice struct {
	dir    *tricorder.DirectorySpec
	size   uint64
	probed bool
}

func Register(dir *tricorder.DirectorySpec) *prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
