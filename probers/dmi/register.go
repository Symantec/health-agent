package dmi

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

const sysfsDir = "/sys/class/dmi/id"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	directory, err := os.Open(sysfsDir)
	if err != nil {
		return p
	}
	defer directory.Close()
	names, err := directory.Readdirnames(-1)
	if err != nil {
		return p
	}
	for _, name := range names {
		switch name {
		case "modalias":
			continue
		case "uevent":
			continue
		}
		pathname := filepath.Join(sysfsDir, name)
		fi, err := os.Lstat(pathname)
		if err != nil {
			continue
		}
		if fi.Mode()&os.ModeType != 0 {
			continue
		}
		if file, err := os.Open(pathname); err != nil {
			continue
		} else {
			buffer := make([]byte, 256)
			if nRead, err := file.Read(buffer); err == nil && nRead > 0 {
				str := strings.TrimSpace(string(buffer[:nRead]))
				if len(str) > 0 {
					dir.RegisterMetric(path.Join("id", name), &str, units.None,
						"")
				}
			}
			file.Close()
		}
	}
	return p
}
