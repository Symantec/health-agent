package url

import (
	"net/http"

	"github.com/Symantec/tricorder/go/tricorder"
)

type urlconfig struct {
	testname            string
	urlpath             string
	urlport             uint
	useTLS              bool
	hasTricorderMetrics bool
	healthy             bool
	statusCode          uint
	httpTransport       *http.Transport
	httpClient          *http.Client
	error               string
}

func Makeurlprober(testname string, urlpath string, urlport uint) *urlconfig {
	return newUrlProber(testname, urlpath, urlport)
}

func (p *urlconfig) GetPort() (uint, bool) {
	return uint(p.urlport), p.useTLS
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
