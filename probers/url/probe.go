package url

import (
	"fmt"
	"net/http"
)

func (p *urlconfig) probe() error {
	address := fmt.Sprintf("http://localhost:%d%s", p.urlport, p.urlpath)
	res, err := http.Get(address)
	if err != nil {
		p.healthy = false
		p.error = fmt.Sprintf("%s", err)
		return err
	}
	defer res.Body.Close()
	p.statusCode = res.StatusCode
	if res.StatusCode == 200 {
		p.healthy = true
		p.error = ""
	} else {
		p.healthy = false
		p.error = "Status not OK"
	}
	return nil
}
