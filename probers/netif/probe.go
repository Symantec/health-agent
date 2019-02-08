package netif

import (
	"bufio"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
	"strings"
	"time"
)

var filename string = "/proc/net/dev"
var sysfsFilenameFormat string = "/sys/class/net/%s/%s"

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
		var err error
		netIf.dir, err = p.dir.RegisterDirectory(netIfName)
		if err != nil {
			return err
		}
		netIf.name = netIfName
		hwAddress, err := readSysfsString(netIfName, "address")
		if err == nil {
			if err := netIf.dir.RegisterMetric("address", &hwAddress,
				units.None, "hardware (MAC) address"); err != nil {
				return err
			}
		}
		netIf.carrier, err = readSysfsBool(netIfName, "carrier")
		if err == nil {
			if err := netIf.dir.RegisterMetric("carrier", &netIf.carrier,
				units.None, "true if carrier detected"); err != nil {
				return err
			}
		}
		mtu, err := readSysfsUint64(netIfName, "mtu")
		if err == nil {
			if err := netIf.dir.RegisterMetric("mtu", &mtu, units.Byte,
				"Messate Transfer Unit"); err != nil {
				return err
			}
		}
		if err := netIf.dir.RegisterMetric("multicast-frames",
			&netIf.multicastFrames, units.None,
			"total multicast frames received or transmitted"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-compressed-packets",
			&netIf.rxCompressedPackets, units.None,
			"compressed packets received"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-data", &netIf.rxData.value,
			units.Byte, "bytes received"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-dropped", &netIf.rxDropped,
			units.None, "receive packets dropped"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-errors", &netIf.rxErrors,
			units.None, "total receive errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-frame-errors",
			&netIf.rxFrameErrors, units.None,
			"receive framing errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-overruns", &netIf.rxOverruns,
			units.None, "receive overrun errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("rx-packets", &netIf.rxPackets.value,
			units.None, "total packets received"); err != nil {
			return err
		}
		netIf.speed, err = readSysfsUint64(netIfName, "speed")
		if err == nil {
			netIf.speed *= 1000000 / 8
			if err := netIf.dir.RegisterMetric("speed", &netIf.speed,
				units.None, "link speed in Bytes/Second"); err != nil {
				return err
			}
		}
		if err := netIf.dir.RegisterMetric("tx-carrier-losses",
			&netIf.txCarrierLosses, units.None,
			"transmit carrier losses"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-collision-errors",
			&netIf.txCollisionErrors, units.None,
			"transmit collision errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-compressed-packets",
			&netIf.txCompressedPackets, units.None,
			"compressed packets transmitted"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-data", &netIf.txData.value,
			units.Byte, "bytes transmitted"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-dropped", &netIf.txDropped,
			units.None, "transmit packets dropped"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-errors", &netIf.txErrors,
			units.None, "total transmit errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-overruns", &netIf.txOverruns,
			units.None, "transmit overrun errors"); err != nil {
			return err
		}
		if err := netIf.dir.RegisterMetric("tx-packets", &netIf.txPackets.value,
			units.None, "total packets transmitted"); err != nil {
			return err
		}
	}
	netIf.probed = true
	nScanned, err := fmt.Sscanf(netIfData,
		"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&netIf.rxData.tmp, &netIf.rxPackets.tmp, &netIf.rxErrors,
		&netIf.rxDropped, &netIf.rxOverruns, &netIf.rxFrameErrors,
		&netIf.rxCompressedPackets,
		&netIf.multicastFrames,
		&netIf.txData.tmp, &netIf.txPackets.tmp, &netIf.rxErrors,
		&netIf.txDropped, &netIf.txOverruns, &netIf.txCollisionErrors,
		&netIf.txCarrierLosses, &netIf.txCompressedPackets)
	if err != nil {
		return err
	}
	if nScanned < 16 {
		return fmt.Errorf("only read %d values from %s", nScanned, line)
	}
	// Handle integer overflow.
	netIf.rxData.update()
	netIf.rxPackets.update()
	netIf.txData.update()
	netIf.txPackets.update()
	netIf.carrier, _ = readSysfsBool(netIfName, "carrier")
	currentTime := time.Now()
	if !netIf.lastProbeTime.IsZero() {
		duration := currentTime.Sub(netIf.lastProbeTime)
		netIf.rxDataRate = uint64(float64(netIf.rxData.value-netIf.lastRxData) /
			duration.Seconds())
		netIf.txDataRate = uint64(float64(netIf.txData.value-netIf.lastTxData) /
			duration.Seconds())
	}
	netIf.lastProbeTime = currentTime
	netIf.lastRxData = netIf.rxData.value
	netIf.lastTxData = netIf.txData.value
	return nil
}

func (wi *wrappingInt) update() {
	if wi.offset+wi.tmp < wi.value {
		wi.offset += 1 << 32
	}
	wi.value = wi.offset + wi.tmp
}

func readSysfsUint64(netIfName, filename string) (uint64, error) {
	filename = fmt.Sprintf(sysfsFilenameFormat, netIfName, filename)
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var value uint64
	nScanned, err := fmt.Fscanf(file, "%d", &value)
	if err != nil {
		return 0, err
	}
	if nScanned < 1 {
		return 0, fmt.Errorf("only read %d values from: %s", nScanned, filename)
	}
	return value, nil
}

func readSysfsBool(netIfName, filename string) (bool, error) {
	filename = fmt.Sprintf(sysfsFilenameFormat, netIfName, filename)
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()
	var ivalue uint
	nScanned, err := fmt.Fscanf(file, "%d", &ivalue)
	if err != nil {
		return false, err
	}
	if nScanned < 1 {
		return false, fmt.Errorf("only read %d values from: %s",
			nScanned, filename)
	}
	if ivalue == 0 {
		return false, nil
	}
	return true, nil
}

func readSysfsString(netIfName, filename string) (string, error) {
	filename = fmt.Sprintf(sysfsFilenameFormat, netIfName, filename)
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var value string
	nScanned, err := fmt.Fscanf(file, "%s", &value)
	if err != nil {
		return "", err
	}
	if nScanned < 1 {
		return "", fmt.Errorf("only read %d values from: %s",
			nScanned, filename)
	}
	return value, nil
}
