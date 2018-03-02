package memory

import (
	"github.com/Symantec/Dominator/lib/meminfo"
)

func (p *prober) probe() error {
	if info, err := meminfo.GetMemInfo(); err != nil {
		return err
	} else {
		p.available = info.Available
		p.free = info.Free
		p.total = info.Total
	}
	return nil
}
