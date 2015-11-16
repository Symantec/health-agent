package netif

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
	"strings"
)

var filename string = "/proc/net/dev"

func (p *prober) probe() error {
	for _, netIf := range p.netInterfaces {
		netIf.probed = false
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := p.processNetdevLine(scanner.Text()); err != nil {
			return err
		}
	}
	// TODO(rgooch): Clean up unprobed network interfaces once tricorder
	//               supports unregistration.
	return scanner.Err()
}

func (p *prober) processNetdevLine(line string) error {
	splitLine := strings.SplitN(line, ":", 2)
	if len(splitLine) != 2 {
		return nil
	}
	netIfName := strings.TrimSpace(splitLine[0])
	netIfData := splitLine[1]
	netIf := p.netInterfaces[netIfName]
	if netIf == nil {
		netIf = new(netInterface)
		p.netInterfaces[netIfName] = netIf
		metricsDir, err := p.dir.RegisterDirectory(netIfName)
		if err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-data", &netIf.rxData,
			units.Byte, "bytes received"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-packets", &netIf.rxPackets,
			units.None, "packets received"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-data", &netIf.txData,
			units.Byte, "bytes transmitted"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-packets", &netIf.txPackets,
			units.None, "packets transmitted"); err != nil {
			return err
		}
	}
	netIf.probed = true
	var ign uint64
	nScanned, err := fmt.Sscanf(netIfData,
		"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&netIf.rxData, &netIf.rxPackets, &ign, &ign, &ign, &ign, &ign, &ign,
		&netIf.txData, &netIf.txPackets, &ign, &ign, &ign, &ign, &ign, &ign)
	if err != nil {
		return err
	}
	if nScanned < 16 {
		return errors.New(fmt.Sprintf("only read %d values from %s",
			nScanned, line))
	}
	return nil
}
