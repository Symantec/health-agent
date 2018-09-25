package storage

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	healthGood     = "good"
	healthMarginal = "marginal"
	healthAtRisk   = "at risk"
	healthFailed   = "failed"
	healthTimedOut = "health check timed out"
	healthUnknown  = "health unknown"

	adapterHeaderSep = "================"
)

type resultsType map[string]map[string]string

func getIntFromMap(table map[string]string, key string) int64 {
	if strValue, ok := table[key]; !ok {
		return -1
	} else {
		value, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			return -1
		}
		return value
	}
}

func megaCliProbe(megaCliPath string) resultsType {
	cmd := exec.Command(megaCliPath, "-AdpAllInfo", "-aALL", "-nolog")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	if err := cmd.Start(); err != nil {
		return nil
	}
	results, err := parseMegaCliAdapter(pipe)
	if err != nil {
		return nil
	}
	if err := cmd.Wait(); err != nil {
		return nil
	}
	return results
}

func parseMegaCliAdapter(r io.Reader) (resultsType, error) {
	raidStats := make(resultsType)
	scanner := bufio.NewScanner(r)
	header := ""
	last := ""
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == adapterHeaderSep {
			header = last
			raidStats[header] = map[string]string{}
			continue
		}
		last = text
		if header == "" { // skip Adapter #X and separator
			continue
		}
		parts := strings.SplitN(text, ":", 2)
		if len(parts) != 2 {
			// This section never includes anything we are interested in.
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		raidStats[header][key] = value
	}
	return raidStats, nil
}

func parseResults(results resultsType) string {
	if results == nil {
		return healthUnknown
	}
	if deviceResults, ok := results["Device Present"]; !ok {
		return healthUnknown
	} else {
		if value := getIntFromMap(deviceResults, "Offline"); value < 0 {
			return healthUnknown
		} else if value > 0 {
			return healthFailed
		}
		if value := getIntFromMap(deviceResults, "Degraded"); value < 0 {
			return healthUnknown
		} else if value > 0 {
			return healthAtRisk
		}
		if value := getIntFromMap(deviceResults, "Critical Disks"); value < 0 {
			return healthUnknown
		} else if value > 0 {
			return healthMarginal
		}
		if value := getIntFromMap(deviceResults, "Failed Disks"); value < 0 {
			return healthUnknown
		} else if value > 0 {
			return healthMarginal
		}
		return healthGood
	}
}

func (p *prober) loopMegaCliProbe(megaCliPath string) {
	for {
		p.megaCliProbe(megaCliPath)
		time.Sleep(time.Minute * 3)
	}
}

func (p *prober) megaCliProbe(megaCliPath string) {
	resultsChan := make(chan resultsType)
	timer := time.NewTimer(time.Minute * 2)
	go func() {
		resultsChan <- megaCliProbe(megaCliPath)
	}()
	var results resultsType
	select {
	case results = <-resultsChan:
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
		p.health = healthTimedOut
		results = <-resultsChan
	}
	p.health = parseResults(results)
}
