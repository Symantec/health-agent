package aws

import (
	"bytes"
	"errors"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"io"
	"net/http"
	"strings"
)

var versionFilename string = "/proc/version"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	instanceId, err := getInstanceId()
	if err == nil {
		if err := dir.RegisterMetric(
			"instance-id",
			&instanceId,
			units.None,
			"AWS instance ID"); err != nil {
			panic(err)
		}
	}
	return p
}

func getInstanceId() (string, error) {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
