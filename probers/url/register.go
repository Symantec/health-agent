package url

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

func newUrlProber(testname string, urlpath string, urlport uint) *urlconfig {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &urlconfig{
		testname:      testname,
		urlpath:       urlpath,
		urlport:       urlport,
		httpTransport: httpTransport,
		httpClient: &http.Client{
			Transport: httpTransport,
			Timeout:   time.Second * 10,
		},
	}
}

func (p *urlconfig) register(dir *tricorder.DirectorySpec) error {
	if err := dir.RegisterMetric("error", &p.error,
		units.None, "Error if any"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("has-tricorder-metrics",
		&p.hasTricorderMetrics, units.None, "Tricorder metrics?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("healthy", &p.healthy,
		units.None, "Healthy?"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("port-number",
		&p.urlport, units.None, "Port number"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("status-code", &p.statusCode,
		units.None, "Status code"); err != nil {
		return err
	}
	if err := dir.RegisterMetric("use-tls",
		&p.useTLS, units.None, "Connect with TLS?"); err != nil {
		return err
	}
	return nil
}
