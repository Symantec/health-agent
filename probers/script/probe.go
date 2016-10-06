package script

import (
	"fmt"
	"os"
	"os/exec"
)

func (p *scriptconfig) probe() error {
	if _, err := os.Stat(p.scriptfilepath); os.IsNotExist(err) {
		p.runSuccessful = false
		p.exitError = fmt.Sprint(err)
		return err
	}
	if err := exec.Command(p.scriptfilepath).Run(); err != nil {
		p.runSuccessful = false
		p.exitError = fmt.Sprint(err)
		return err
	}
	p.runSuccessful = true
	p.exitError = ""
	return nil
}
