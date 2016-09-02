package kernel

import (
	"bufio"
	"fmt"
	"os"
)

var randomEntropyFilename = "/proc/sys/kernel/random/entropy_avail"

func (p *prober) probe() error {
	file, err := os.Open(randomEntropyFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	nScanned, err := fmt.Fscanf(reader, "%d", &p.randomEntropyAvailable)
	if nScanned != 1 {
		return err
	}
	return nil
}
