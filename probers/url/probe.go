package url

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (p *urlconfig) probe() error {
	urlpath := p.urlpath
	if !strings.HasPrefix(urlpath, "http://") {
		urlpath = "http://" + urlpath
	}
	urlpath += ":" + strconv.Itoa(p.urlport)
	res, err := http.Get(urlpath)
	if err != nil {
		p.healthy = false
		p.error = fmt.Sprintf("%s", err)
		return err
	}
	p.statusCode = res.StatusCode
	if strings.Contains(res.Status, "OK") {
		p.healthy = true
		p.error = ""
	} else {
		p.healthy = false
		p.error = "Status not OK"
	}
	return nil
}
