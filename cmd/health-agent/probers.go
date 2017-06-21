package main

import (
	"github.com/Symantec/health-agent/lib/proberlist"
	"github.com/Symantec/health-agent/probers/aws"
	"github.com/Symantec/health-agent/probers/filesystems"
	"github.com/Symantec/health-agent/probers/kernel"
	"github.com/Symantec/health-agent/probers/memory"
	"github.com/Symantec/health-agent/probers/netif"
	"github.com/Symantec/health-agent/probers/network"
	"github.com/Symantec/health-agent/probers/packages"
	"github.com/Symantec/health-agent/probers/scheduler"
	"github.com/Symantec/health-agent/probers/storage"
	"github.com/Symantec/health-agent/probers/systime"
)

func setupProbers() (*proberlist.ProberList, error) {
	pl := proberlist.New("/probers")
	go func() { pl.Add(aws.New(), "/sys/cloud/aws", 0) }()
	pl.CreateAndAdd(filesystems.Register, "/sys/fs", 0)
	pl.CreateAndAdd(scheduler.Register, "/sys/sched", 0)
	pl.CreateAndAdd(memory.Register, "/sys/memory", 0)
	pl.CreateAndAdd(netif.Register, "/sys/netif", 0)
	pl.CreateAndAdd(network.Register, "/sys/network", 0)
	pl.CreateAndAdd(storage.Register, "/sys/storage", 0)
	pl.CreateAndAdd(systime.Register, "/sys/systime", 0)
	pl.CreateAndAdd(kernel.Register, "/sys/kernel", 0)
	pl.CreateAndAdd(packages.Register, "/sys/packages", 0)
	return pl, nil
}
