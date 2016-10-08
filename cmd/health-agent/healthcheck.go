package main

import (
	libprober "github.com/Symantec/health-agent/lib/proberlist"
	pidprober "github.com/Symantec/health-agent/probers/pidfile"
	testprogprober "github.com/Symantec/health-agent/probers/testprog"
	urlprober "github.com/Symantec/health-agent/probers/url"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type testConfig struct {
	Testtype  string `yaml:"type"`
	Probefreq uint8  `yaml:"probe-freq"`
	Specs     testSpecs
}

type testSpecs struct {
	Pathname string
	Urlpath  string `yaml:"url-path"`
	Urlport  int    `yaml:"url-port"`
}

func setupHealthchecks(configDir string, pl *libprober.ProberList,
	logger *log.Logger) error {
	topDir := "/health-checks"
	configdir, err := os.Open(path.Join(configDir, "tests.d"))
	defer configdir.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	configfiles, err := configdir.Readdir(0)
	if err != nil {
		return err
	}
	for _, configfile := range configfiles {
		if configfile.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(path.Join(configdir.Name(),
			configfile.Name()))
		if err != nil {
			logger.Printf("Unable to read file %q: %s",
				configfile.Name(), err)
			return err
		}
		c := testConfig{}
		if err := yaml.Unmarshal([]byte(data), &c); err != nil {
			logger.Printf("Error unmarshalling file %s: %q",
				configfile.Name(), err)
			return err
		}
		testname := strings.Split(configfile.Name(), ".")[0]
		if prober := makeProber(testname, &c, logger); prober != nil {
			pl.Add(prober, path.Join(topDir, testname), c.Probefreq)
		}
	}
	return nil
}

func makeProber(testname string, c *testConfig,
	logger *log.Logger) libprober.RegisterProber {
	switch c.Testtype {
	case "pid":
		pidpath := c.Specs.Pathname
		if pidpath == "" {
			return nil
		}
		return pidprober.Makepidprober(testname, pidpath)
	case "testprog":
		testprogpath := c.Specs.Pathname
		if testprogpath == "" {
			return nil
		}
		return testprogprober.Maketestprogprober(testname, testprogpath)
	case "url":
		urlpath := c.Specs.Urlpath
		urlport := c.Specs.Urlport
		if urlpath == "" {
			return nil
		}
		return urlprober.Makeurlprober(testname, urlpath, urlport)
	default:
		logger.Printf("Test type %s not supported", c.Testtype)
		return nil
	}
}
