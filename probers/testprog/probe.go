package testprog

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func (p *testprogconfig) probe() error {
	if _, err := os.Stat(p.filepath); os.IsNotExist(err) {
		p.sethealthy(false, err)
		return err
	}
	cmd := exec.Command(p.filepath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		p.sethealthy(false, err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		p.sethealthy(false, err)
		return err
	}
	if err := cmd.Start(); err != nil {
		p.sethealthy(false, err)
		return err
	}
	var buffout bytes.Buffer
	var buffer bytes.Buffer
	scanout := bufio.NewScanner(stdout)
	for scanout.Scan() {
		buffout.WriteString(scanout.Text())
	}
	if err := scanout.Err(); err != nil {
		p.sethealthy(false, err)
		return err
	}
	p.stdout = buffout.String()
	scanerr := bufio.NewScanner(stderr)
	for scanerr.Scan() {
		buffer.WriteString(scanerr.Text())
	}
	if err := scanerr.Err(); err != nil {
		p.sethealthy(false, err)
		return err
	}
	p.stderr = buffer.String()
	if err := cmd.Wait(); err != nil {
		p.sethealthy(false, err)
		return err
	}
	p.sethealthy(true, nil)
	p.exitCode = 0
	return nil
}

func (p *testprogconfig) sethealthy(val bool, err error) {
	p.healthy = val
	p.exitError = fmt.Sprint(err)
}
