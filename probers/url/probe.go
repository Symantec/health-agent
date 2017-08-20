package url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const hasTricorderUrl = "/has-tricorder-metrics"

func (p *urlconfig) probe() error {
	address := fmt.Sprintf("http://localhost:%d%s", p.urlport, p.urlpath)
	res, err := http.Get(address)
	if err != nil {
		p.healthy = false
		p.error = err.Error()
		return err
	}
	defer res.Body.Close()
	p.statusCode = uint(res.StatusCode)
	if res.StatusCode == 200 {
		p.healthy = true
		p.error = ""
		if hasTricorderMetrics, err := p.probeTricorder(); err == nil {
			p.hasTricorderMetrics = hasTricorderMetrics
		}
	} else {
		p.healthy = false
		p.error = res.Status
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			status := strings.TrimSpace(string(body))
			if status != "" {
				p.error = status
			}
		}
	}
	return nil
}

func (p *urlconfig) probeTricorder() (bool, error) {
	address := fmt.Sprintf("http://localhost:%d%s", p.urlport, hasTricorderUrl)
	res, err := http.Get(address)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}
