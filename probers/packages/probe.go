package packages

import (
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *prober) probe() error {
	if err := p.probeDebs(); err != nil {
		return err
	}
	if err := p.probeRpms(); err != nil {
		return err
	}
	return nil
}

func (p *prober) addPackage(pList *packageList, pEntry *packageEntry) error {
	dir, err := p.debian.dir.RegisterDirectory(pEntry.name)
	if err != nil {
		return err
	}
	pEntry.dir = dir
	if err := dir.RegisterMetric("size", &pEntry.size, units.Byte,
		"package size"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("version", &pEntry.version, units.None,
		"package version"); err != nil {
		return err
	}
	pList.packages[pEntry.name] = pEntry
	return nil
}
