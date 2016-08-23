package packages

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func probeRpms(pList *packageList, reader io.Reader) error {
	cmd := exec.Command("rpm", "-qa", "--queryformat",
		"%{NAME} %{VERSION}-%{RELEASE} %{SIZE}\n")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
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
		if len(fields) != 3 {
			continue
		}
		pEntry := &packageEntry{name: fields[0], version: fields[1]}
		var err error
		if pEntry.size, err = strconv.ParseUint(fields[2], 10, 64); err != nil {
			return err
		}
		if err := addPackage(pList, pEntry); err != nil {
			return err
		}
	}
	return nil
}
