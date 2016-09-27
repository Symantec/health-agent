package proberlist

import (
	"github.com/Symantec/health-agent/lib/prober"
	"github.com/Symantec/tricorder/go/tricorder"
	"io"
	"log"
	"time"
)

// HtmlWriter defines a type that can write a HTML snippet about itself.
type HtmlWriter interface {
	WriteHtml(writer io.Writer)
}

type RegisterFunc func(dir *tricorder.DirectorySpec) prober.Prober

type proberType struct {
	prober                prober.Prober
	probeTimeDistribution *tricorder.CumulativeDistribution
}

// ProberList defines a type which manages a list of Probers.
type ProberList struct {
	probers               []proberType
	proberPath            string
	probeStartTime        time.Time
	probeTimeDistribution *tricorder.CumulativeDistribution
}

// New returns a new ProberList. Only one should be created per application.
// Metrics showing the operation of the Probers (not the metrics that the
// Probers collect) will be placed under proberPath.
func New(proberPath string) *ProberList {
	return newProberList(proberPath)
}

// CreateAndAdd registers a new Prober which is created by the registerFunc. The
// path for the metrics for the Prober is given by path. The preferred probe
// interval in seconds is given by probeInterval.
func (pl *ProberList) CreateAndAdd(registerFunc RegisterFunc, path string,
	probeInterval uint8) {
	pl.createAndAdd(registerFunc, path, probeInterval)
}

// StartProbing creates one or more goroutines which will run probes in an
// infinite loop. The default probe interval in seconds is given by
// defaultProbeInterval. The logger will be used to log problems.
func (pl *ProberList) StartProbing(defaultProbeInterval uint,
	logger *log.Logger) {
	go pl.proberLoop(defaultProbeInterval, logger)
}

// WriteHtml will write HTML snippets to writer. Each Prober that implements the
// HtmlWriter interface will have it's WriteHtml method called. These methods
// are called in the order in which the Probers were added.
func (pl *ProberList) WriteHtml(writer io.Writer) {
	pl.writeHtml(writer)
}
