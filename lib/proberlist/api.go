package proberlist

import (
	"github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
	"io"
	"log"
	"time"
)

type RegisterFunc func(dir *tricorder.DirectorySpec) prober.Prober

type proberType struct {
	prober                prober.Prober
	probeTimeDistribution *tricorder.CumulativeDistribution
}

type ProberList struct {
	probers               []proberType
	proberPath            string
	probeStartTime        time.Time
	probeTimeDistribution *tricorder.CumulativeDistribution
}

func New(proberPath string) *ProberList {
	return newProberList(proberPath)
}

func (pl *ProberList) Add(registerFunc RegisterFunc, path string,
	probeInterval uint8) {
	pl.add(registerFunc, path, probeInterval)
}

func (pl *ProberList) StartProbing(defaultProbeInterval uint,
	logger *log.Logger) {
	go pl.proberLoop(defaultProbeInterval, logger)
}

func (pl *ProberList) WriteHtml(writer io.Writer) {
	pl.writeHtml(writer)
}
