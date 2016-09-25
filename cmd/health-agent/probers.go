package main

import (
	"github.com/Symantec/health-agent/lib/proberlist"
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
	pl.Add(filesystems.Register, "/sys/fs", 0)
	pl.Add(scheduler.Register, "/sys/sched", 0)
	pl.Add(memory.Register, "/sys/memory", 0)
	pl.Add(netif.Register, "/sys/netif", 0)
	pl.Add(network.Register, "/sys/network", 0)
	pl.Add(storage.Register, "/sys/storage", 0)
	pl.Add(systime.Register, "/sys/systime", 0)
	pl.Add(kernel.Register, "/sys/kernel", 0)
	pl.Add(packages.Register, "/sys/packages", 0)
	return pl, nil
}
