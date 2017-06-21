package aws

import (
	"github.com/Symantec/health-agent/lib/proberlist"
)

func New() proberlist.RegisterProber {
	return _new()
}
