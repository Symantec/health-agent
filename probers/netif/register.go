package netif

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	p.dir = dir
	p.netInterfaces = make(map[string]*netInterface)
	// TODO(rgooch): Remove this call once tricorder supports dynamic
	//               registration.
	p.probe()
	return p
}
