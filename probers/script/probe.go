package script

import (
	"fmt"
	"os"
	"os/exec"
)

func (p *scriptconfig) probe() error {
	if _, err := os.Stat(p.scriptfilepath); os.IsNotExist(err) {
		p.healthy = false
		p.exitError = fmt.Sprint(err)
		return err
	}
	if err := exec.Command(p.scriptfilepath).Run(); err != nil {
		p.healthy = false
		p.exitError = fmt.Sprint(err)
		return err
	}
	p.healthy = true
	p.exitError = ""
	return nil
}
