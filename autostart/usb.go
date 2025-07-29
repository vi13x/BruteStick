package autostart

import (
	"os"
	"strings"
)

func IsRunningFromUSB() bool {
	exePath, _ := os.Executable()
	for _, drive := range []string{"D:", "E:", "F:", "G:"} {
		if strings.HasPrefix(strings.ToUpper(exePath), drive) {
			return true
		}
	}
	return false
}
