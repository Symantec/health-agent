package memory

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var filename string = "/proc/meminfo"

func (p *prober) probe() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := p.processMeminfoLine(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (p *prober) processMeminfoLine(line string) error {
	splitLine := strings.SplitN(line, ":", 2)
	if len(splitLine) != 2 {
		return nil
	}
	meminfoName := splitLine[0]
	meminfoDataString := strings.TrimSpace(splitLine[1])
	var ptr *uint64
	switch meminfoName {
	case "MemAvailable":
		ptr = &p.available
	case "MemFree":
		ptr = &p.free
	case "MemTotal":
		ptr = &p.total
	default:
		return nil
	}
	var meminfoData uint64
	var meminfoUnit string
	fmt.Sscanf(meminfoDataString, "%d %s", &meminfoData, &meminfoUnit)
	if meminfoUnit != "kB" {
		return errors.New(fmt.Sprintf("unknown unit: %s for: %s",
			meminfoUnit, meminfoName))
	}
	*ptr = meminfoData * 1024
	return nil
}
