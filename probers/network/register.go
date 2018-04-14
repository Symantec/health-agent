package network

import (
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
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
	resolverDir, err := dir.RegisterDirectory("dns-resolver")
	if err != nil {
		panic(err)
	}
	if err := resolverDir.RegisterMetric("domain", &p.resolverDomain,
		units.None, "default domain name"); err != nil {
		panic(err)
	}
	p.resolverNameservers = tricorder.NewList([]string{},
		tricorder.ImmutableSlice)
	if err := resolverDir.RegisterMetric("nameservers", p.resolverNameservers,
		units.None, "resolvers"); err != nil {
		panic(err)
	}
	p.resolverSearchDomains = tricorder.NewList([]string{},
		tricorder.ImmutableSlice)
	if err := resolverDir.RegisterMetric("search-domains",
		p.resolverSearchDomains, units.None,
		"domains (zones) to search"); err != nil {
		panic(err)
	}
	return p
}
