package url

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type urlconfig struct {
	testname            string
	urlpath             string
	urlport             int
	hasTricorderMetrics bool
	healthy             bool
	statusCode          int
	error               string
}

func Makeurlprober(testname string, urlpath string, urlport int) *urlconfig {
	p := new(urlconfig)
	p.testname = testname
	p.urlpath = urlpath
	p.urlport = urlport
	return p
}

func (p *urlconfig) GetPort() uint {
	return uint(p.urlport)
}

func (p *urlconfig) HasTricorderMetrics() bool {
	return p.hasTricorderMetrics
}

func (p *urlconfig) HealthCheck() bool {
	return p.healthy
}

func (p *urlconfig) Probe() error {
	return p.probe()
}

func (p *urlconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}
