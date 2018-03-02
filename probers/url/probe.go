package url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const hasTricorderUrl = "/has-tricorder-metrics"

var (
	maybeHTTP  = "server gave HTTP response to HTTPS client"
	maybeHTTPS = "malformed HTTP response"
)

func (p *urlconfig) getURL(path string) (*http.Response, error) {
	var url string
	if p.useTLS {
		url = fmt.Sprintf("https://localhost:%d%s", p.urlport, path)
	} else {
		url = fmt.Sprintf("http://localhost:%d%s", p.urlport, path)
	}
	return p.httpClient.Get(url)
}

func (p *urlconfig) probe() error {
	defer p.httpTransport.CloseIdleConnections()
	res, err := p.getURL(p.urlpath)
	if err != nil {
		if p.useTLS && strings.Contains(err.Error(), maybeHTTP) {
			p.useTLS = false
			res, err = p.getURL(p.urlpath)
		} else if !p.useTLS && strings.Contains(err.Error(), maybeHTTPS) {
			p.useTLS = true
			res, err = p.getURL(p.urlpath)
		}
	}
	if err != nil {
		p.healthy = false
		p.error = err.Error()
		return err
	}
	body, bodyErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
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
		if bodyErr == nil {
			status := strings.TrimSpace(string(body))
			if status != "" {
				p.error = status
			}
		}
	}
	return nil
}

func (p *urlconfig) probeTricorder() (bool, error) {
	res, err := p.getURL(hasTricorderUrl)
	if err != nil {
		return false, err
	}
	ioutil.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}
