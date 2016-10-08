package url

import (
	"github.com/Symantec/tricorder/go/tricorder"
)

type urlconfig struct {
	testname   string
	urlpath    string
	urlport    int
	healthy    bool
	statusCode int
	error      string
}

func (p *urlconfig) Register(dir *tricorder.DirectorySpec) error {
	return p.register(dir)
}

func (p *urlconfig) Probe() error {
	return p.probe()
}

func Makeurlprober(testname string, urlpath string, urlport int) *urlconfig {
	p := new(urlconfig)
	p.testname = testname
	p.urlpath = urlpath
	p.urlport = urlport
	return p
}
