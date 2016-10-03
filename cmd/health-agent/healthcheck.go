package main

import (
	libprober "github.com/Symantec/health-agent/lib/proberlist"
	pidprober "github.com/Symantec/health-agent/probers/pidfile"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type testConfig struct {
	Testtype  string `yaml:"type"`
	Probefreq uint8 `yaml:omitempty`
	Specs     testSpecs
}

type testSpecs struct {
	Pidfilepath string
	Urlpath string
	Urlport uint
}

func setupHealthchecks(configDir string, pl *libprober.ProberList, logger *log.Logger) error {
	topDir := "/health-checks/"
	configdir, err := os.Open(configDir)
	if err != nil {
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
		data, err := ioutil.ReadFile(configDir + configfile.Name())
		if err != nil {
			logger.Printf("Unable to read file %s\n%q", configfile.Name(), err)
			continue
		}
		c := testConfig{}
		if err := yaml.Unmarshal([]byte(data), &c); err != nil {
			logger.Printf("Error unmarshalling file %s: %q", configfile.Name(), err)
			continue
		}
		testname := strings.Split(configfile.Name(), ".")[0]
		prober := makeProber(testname, &c, logger)
		if prober != nil {
			pl.Add(prober, topDir + testname, c.Probefreq)
		}
	}
	return nil
}

func makeProber(testname string, c *testConfig, logger *log.Logger) (libprober.RegisterProber) {
	switch c.Testtype {
	case "pid":
		pidpath := c.Specs.Pidfilepath
		if pidpath == "" {
			return nil 
		}
		return pidprober.Makepidprober(testname, pidpath)
	default:
		logger.Println("Test type %s not supported", c.Testtype)
		return nil
	}
}
