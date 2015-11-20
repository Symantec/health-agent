package main

import (
	"github.com/Symantec/health-agent/probers/filesystems"
	"github.com/Symantec/health-agent/probers/loadavg"
	"github.com/Symantec/health-agent/probers/memory"
	"github.com/Symantec/health-agent/probers/netif"
	"github.com/Symantec/health-agent/probers/network"
	"github.com/Symantec/health-agent/probers/storage"
	"github.com/Symantec/health-agent/probers/systime"
	"github.com/Symantec/tricorder/go/tricorder"
)


type prober interface {
	Probe() error
}

func setupProbers() ([]prober, error) {
	topMetricsDir, err := tricorder.RegisterDirectory("/sys")
	if err != nil {
		return nil, err
	}
	var probers []prober
	probers = append(probers, filesystems.Register(mkdir(topMetricsDir, "fs")))
	probers = append(probers, loadavg.Register(mkdir(topMetricsDir, "loadavg")))
	probers = append(probers, memory.Register(mkdir(topMetricsDir, "memory")))
	probers = append(probers, netif.Register(mkdir(topMetricsDir, "netif")))
	probers = append(probers, network.Register(mkdir(topMetricsDir, "network")))
	probers = append(probers, storage.Register(mkdir(topMetricsDir, "storage")))
	probers = append(probers, systime.Register(mkdir(topMetricsDir, "")))
	return probers, nil
}

func mkdir(dir *tricorder.DirectorySpec, name string) *tricorder.DirectorySpec {
	if name == "" {
		return dir
	}
	subdir, err := dir.RegisterDirectory(name)
	if err != nil {
		panic(err)
	}
	return subdir
}
