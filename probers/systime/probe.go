package systime

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var filename string = "/proc/uptime"

func (p *prober) probe() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var idleTime, upTime float64
	nScanned, err := fmt.Fscanf(file, "%f %f", &upTime, &idleTime)
	if err != nil {
		return err
	}
	if nScanned < 2 {
		return errors.New(fmt.Sprintf("only read %d values from %s",
			nScanned, filename))
	}
	p.idleTime = time.Duration(idleTime * float64(time.Second))
	p.probeTime = time.Now()
	p.upTime = time.Duration(upTime * float64(time.Second))
	return nil
}
