//go:build windows
// +build windows

package registry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

func SetupAutostart(exePath string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("cannot open registry key: %v", err)
	}
	defer k.Close()

	v, _, err := k.GetStringValue("BruteStick")
	if err == nil && v == exePath {
		// Уже настроен автозапуск
		return nil
	}

	return k.SetStringValue("BruteStick", exePath)
}
