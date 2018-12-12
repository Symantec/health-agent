package filesystems

import (
	"bufio"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
	"strings"
	"syscall"
)

var filename string = "/proc/mounts"

func (p *prober) probe() error {
	for _, fs := range p.fileSystems {
		fs.probed = false
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := p.processMountLine(scanner.Text()); err != nil {
			return err
		}
	}
	// TODO(rgooch): Clean up unprobed file-systems once tricorder supports
	//               unregistration.
	return scanner.Err()
}

func (p *prober) processMountLine(line string) error {
	var device, mountPoint, fsType, fsOptions, junk string
	nScanned, err := fmt.Sscanf(line, "%s %s %s %s %s",
		&device, &mountPoint, &fsType, &fsOptions, &junk)
	if err != nil {
		return err
	}
	if nScanned < 4 {
		return fmt.Errorf("only read %d values from %s", nScanned, line)
	}
	if !strings.HasPrefix(device, "/dev/") {
		return nil // Not a device: ignore.
	}
	device = device[5:]
	var statbuf syscall.Statfs_t
	if fd, err := syscall.Open(mountPoint, syscall.O_RDONLY, 0); err != nil {
		return fmt.Errorf("error opening: %s %s", mountPoint, err)
	} else {
		defer syscall.Close(fd)
		if err := syscall.Fstatfs(fd, &statbuf); err != nil {
			return nil // Something weird like FUSE: ignore.
		}
	}
	fs := p.fileSystems[device]
	if fs == nil {
		fs = new(fileSystem)
		p.fileSystems[device] = fs
		if mountPoint == "/" {
			fs.dir = p.dir
		} else {
			fs.dir, err = p.dir.RegisterDirectory(mountPoint)
			if err != nil {
				return err
			}
		}
		fs.mountPoint = mountPoint
		metricsDir, err := fs.dir.RegisterDirectory("METRICS")
		if err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("available", &fs.available,
			units.Byte, "space available to unprivileged users"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("device", &fs.device,
			units.None, "device "); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("free", &fs.free,
			units.Byte, "free space"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("options", &fs.options,
			units.None, "options "); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("size", &fs.size,
			units.Byte, "available space"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("writable", &fs.writable,
			units.None, "true if writable, else read-only"); err != nil {
			return err
		}
	}
	if fs.probed {
		return nil
	}
	fs.available = statbuf.Bavail * uint64(statbuf.Bsize)
	fs.free = statbuf.Bfree * uint64(statbuf.Bsize)
	fs.size = statbuf.Blocks * uint64(statbuf.Bsize)
	fs.device = device
	fs.options = fsOptions
	fs.writable = strings.HasPrefix(fs.options, "rw,")
	fs.probed = true
	return nil
}
