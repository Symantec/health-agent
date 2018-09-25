package storage

import (
	libprober "github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
)

type prober struct {
	dir            *tricorder.DirectorySpec
	health         string
	storageDevices map[string]*storageDevice // map[name]*storageDevice
}

type storageDevice struct {
	dir    *tricorder.DirectorySpec
	size   uint64
	probed bool
}

func Register(dir *tricorder.DirectorySpec) libprober.Prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}
