package network

import (
	"bufio"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"net"
	"os"
)

var filename string = "/proc/net/route"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := p.findGateway(); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("gateway-address", &p.gatewayAddress,
		units.None, "gateway address"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("gateway-interface", &p.gatewayInterfaceName,
		units.None, "gateway interface"); err != nil {
		panic(err)
	}
	latencyBucketer := tricorder.NewGeometricBucketer(0.1, 10e3)
	p.gatewayPingTimeDistribution = latencyBucketer.NewCumulativeDistribution()
	if err := dir.RegisterMetric("gateway-ping-time",
		p.gatewayPingTimeDistribution,
		units.Millisecond, "ping time to gateway"); err != nil {
		panic(err)
	}
	p.gatewayRttDistribution = latencyBucketer.NewCumulativeDistribution()
	if err := dir.RegisterMetric("gateway-rtt", p.gatewayRttDistribution,
		units.Millisecond, "round-trip time to gateway"); err != nil {
		panic(err)
	}
	return p
}

func (p *prober) findGateway() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var interfaceName string
		var destAddr, gatewayAddr, flags, mask uint32
		var ign int
		nCopied, err := fmt.Sscanf(scanner.Text(),
			"%s %x %x %x %d %d %d %x %d %d %d",
			&interfaceName, &destAddr, &gatewayAddr, &flags, &ign, &ign, &ign,
			&mask, &ign, &ign, &ign)
		if err != nil || nCopied < 11 {
			continue
		}
		if destAddr == 0 && flags&0x2 == 0x2 && flags&0x1 == 0x1 {
			p.gatewayAddress = intToIP(gatewayAddr).String()
			p.gatewayInterfaceName = interfaceName
			return nil
		}
	}
	return scanner.Err()
}

func intToIP(ip uint32) net.IP {
	result := make(net.IP, 4)
	for i, _ := range result {
		result[i] = byte(ip >> uint(8*i))
	}
	return result
}
