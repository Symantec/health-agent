package packages

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := &prober{
		dir: dir,
	}
	return p
}
