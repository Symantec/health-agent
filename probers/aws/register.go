package aws

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"io"
	"net/http"
	"strings"
)

var versionFilename string = "/proc/version"

var (
	metaDataKeyToDescription = map[string]string{
		"ami-id":          "AMI ID",
		"hostname":        "local hostname",
		"instance-id":     "AWS instance ID",
		"instance-type":   "AWS instance type",
		"local-ipv4":      "local ipv4",
		"public-hostname": "public hostname",
		"public-ipv4":     "public ipv4",
	}
)

var (
	kErrNotAws = errors.New("aws: Response not AWS metadata")
)

func isAwsMetaData(response *http.Response) bool {
	return response.Header.Get("Server") == "EC2ws"
}

func register(dir *tricorder.DirectorySpec) *prober {
	var aMetricAdded bool
	if keys, err := getMetaDataKeys(); err == nil {
		for _, key := range keys {
			if desc, ok := metaDataKeyToDescription[key]; ok {
				if value, err := getByMetaDataKey(key); err == nil {
					if err := dir.RegisterMetric(
						key,
						&value,
						units.None,
						desc); err != nil {
						panic(err)
					} else {
						aMetricAdded = true
					}
				}
			}
		}
	}
	// If we couldn't register any metrics, we have no prober
	// this prober.
	if !aMetricAdded {
		return nil
	}
	return new(prober)
}

func getMetaDataKeys() ([]string, error) {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	if !isAwsMetaData(resp) {
		return nil, kErrNotAws
	}
	var result []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func getByMetaDataKey(key string) (string, error) {
	resp, err := http.Get(
		fmt.Sprintf(
			"http://169.254.169.254/latest/meta-data/%s", key))

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	if !isAwsMetaData(resp) {
		return "", kErrNotAws
	}
	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
