package packages

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var (
	debStatusFile = "/var/lib/dpkg/status"
)

func (p *prober) probeDebs() error {
	file, err := os.Open(debStatusFile)
	if err != nil {
		if os.IsNotExist(err) {
			if p.debian != nil {
				p.debian.dir.UnregisterDirectory()
				p.debian = nil
			}
		}
		return err
	}
	defer file.Close()
	if p.debian == nil {
		dir, err := p.dir.RegisterDirectory("debs")
		if err != nil {
			return err
		}
		p.debian = &packageList{
			dir:      dir,
			packages: make(map[string]*packageEntry),
		}
	}
	packagesToDelete := make(map[string]struct{})
	for packageName := range p.debian.packages {
		packagesToDelete[packageName] = struct{}{}
	}
	scanner := bufio.NewScanner(file)
	var pEntry *packageEntry
	for scanner.Scan() {
		line := scanner.Text()
		lineLength := len(line)
		if lineLength < 2 {
			continue
		}
		if line[lineLength-1] == '\n' {
			line = line[:lineLength-1]
		}
		fields := strings.Fields(line)
		if fields[0] == "Package:" {
			if pEntry != nil {
				if err := p.addPackage(p.debian, pEntry); err != nil {
					return err
				}
			}
			pEntry = &packageEntry{name: fields[1]}
			delete(packagesToDelete, pEntry.name)
			continue
		}
		if pEntry == nil {
			continue
		}
		switch fields[0] {
		case "Status:":
			if !strings.Contains(line[8:], "installed") {
				pEntry = nil
				continue
			}
		case "Version:":
			pEntry.version = fields[1]
		case "Installed-Size:":
			if size, err := strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			} else {
				pEntry.size = size << 10
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if pEntry != nil {
		if err := p.addPackage(p.debian, pEntry); err != nil {
			return err
		}
	}
	for packageName := range packagesToDelete {
		p.debian.packages[packageName].dir.UnregisterDirectory()
		delete(p.debian.packages, packageName)
	}
	return nil
}
