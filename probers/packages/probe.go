package packages

import (
	"errors"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"io"
	"os"
	"time"
)

type packagerType struct {
	filename string
	name     string
	prober   func(*packageList, io.Reader) error
}

var packagers = []packagerType{
	{"/var/lib/dpkg/status", "debs", probeDebs},
	{"/var/lib/rpm/Packages", "rpms", probeRpms},
}

func (p *prober) probe() error {
	probeStartTime := time.Now()
	for _, packager := range packagers {
		if err := p.probePackager(packager); err != nil {
			return err
		}
	}
	p.lastProbeStartTime = probeStartTime
	return nil
}

func (p *prober) probePackager(packager packagerType) error {
	pList := p.packagers[packager.name]
	file, err := os.Open(packager.filename)
	if err != nil {
		if os.IsNotExist(err) {
			if pList != nil {
				pList.dir.UnregisterDirectory()
				delete(p.packagers, packager.name)
			}
			return nil
		}
		return err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if fi.ModTime().Before(p.lastProbeStartTime) && pList != nil {
		return nil
	}
	if pList == nil {
		dir, err := p.dir.RegisterDirectory(packager.name)
		if err != nil {
			return errors.New(err.Error() + ": " + packager.name)
		}
		pList = &packageList{
			dir:      dir,
			packages: make(map[string]*packageEntry),
		}
		p.packagers[packager.name] = pList
	}
	pList.packagesAddedDuringProbe = make(map[string]struct{})
	if err := packager.prober(pList, file); err != nil {
		return err
	}
	for packageName := range pList.packages {
		if _, ok := pList.packagesAddedDuringProbe[packageName]; !ok {
			pList.packages[packageName].dir.UnregisterDirectory()
			delete(pList.packages, packageName)
		}
	}
	return nil
}

func addPackage(pList *packageList, pEntry *packageEntry) error {
	pList.packagesAddedDuringProbe[pEntry.name] = struct{}{}
	if oldEntry := pList.packages[pEntry.name]; oldEntry != nil {
		pEntry.dir = oldEntry.dir
		*oldEntry = *pEntry
		return nil
	}
	dir, err := pList.dir.RegisterDirectory(pEntry.name)
	if err != nil {
		if err != nil {
			return errors.New(err.Error() + ": " + pEntry.name)
		}
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
