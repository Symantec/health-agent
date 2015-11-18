package loadavg

import (
	"errors"
	"fmt"
	"os"
)

var filename string = "/proc/loadavg"

func (p *prober) probe() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var str string
	nScanned, err := fmt.Fscanf(file, "%f %f %f %s",
		&p.oneMinute, &p.fiveMinutes, &p.fifteenMinutes, &str)
	if err != nil {
		return err
	}
	if nScanned < 3 {
		return errors.New(fmt.Sprintf("only read %d values from %s",
			nScanned, filename))
	}
	return nil
}
