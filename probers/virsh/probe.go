package virsh

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/Symantec/tricorder/go/tricorder/units"
)

func (p *prober) listDomains() error {
	cmd := exec.Command("virsh", "list", "--all")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 3 {
		return errors.New("insufficient lines")
	}
	if lines[1][0] != '-' {
		return errors.New("missing separator")
	}
	domains := make(map[string]string, len(lines)-2)
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		if len(fields) < 3 {
			return errors.New("bad line: \"" + line + "\"")
		}
		domains[fields[1]] = strings.Join(fields[2:], " ")
	}
	p.listResults = domains
	return nil
}

func (p *prober) probe() error {
	err := p.listDomains()
	if err != nil {
		return err
	}
	// First unregister domains no longer found.
	for name, domain := range p.domains {
		if _, ok := p.listResults[name]; !ok {
			domain.dir.UnregisterDirectory()
			delete(p.domains, name)
		}
	}
	for name, state := range p.listResults {
		info, ok := p.domains[name]
		if ok {
			info.state = state
		} else {
			info = &domainInfo{state: state}
			dir, err := p.domainsDir.RegisterDirectory(name)
			if err != nil {
				return err
			}
			info.dir = dir
			err = dir.RegisterMetric("state", &info.state, units.None,
				"state of virtual machine")
			if err != nil {
				return err
			}
			p.domains[name] = info
		}
	}
	return nil
}
