package packages

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func probeDebs(pList *packageList, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	var pEntry *packageEntry
	for scanner.Scan() {
		line := scanner.Text()
		lineLength := len(line)
		if lineLength < 2 {
			continue
		}
		if line[lineLength-1] == '\n' {
			line = line[:lineLength-1]
		}
		fields := strings.Fields(line)
		if fields[0] == "Package:" {
			if pEntry != nil {
				if err := addPackage(pList, pEntry); err != nil {
					return err
				}
			}
			pEntry = &packageEntry{name: fields[1]}
			continue
		}
		if pEntry == nil {
			continue
		}
		switch fields[0] {
		case "Status:":
			if !strings.Contains(line[8:], "installed") {
				pEntry = nil
				continue
			}
		case "Version:":
			pEntry.version = fields[1]
		case "Installed-Size:":
			if size, err := strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			} else {
				pEntry.size = size << 10
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if pEntry != nil {
		if err := addPackage(pList, pEntry); err != nil {
			return err
		}
	}
	return nil
}
