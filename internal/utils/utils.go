package utils

import (
	"encoding/gob"
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

func DefaultCharSet() []rune {
	return []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func PadPassword(runes []rune, length int, pad rune) []rune {
	for len(runes) < length {
		runes = append(runes, pad)
	}
	return runes
}

func NextPassword(current []rune, charset []rune) ([]rune, bool) {
	for i := len(current) - 1; i >= 0; i-- {
		index := indexOf(charset, current[i])
		if index == -1 {
			current[i] = charset[0]
			continue
		}
		if index+1 < len(charset) {
			current[i] = charset[index+1]
			return current, true
		}
		current[i] = charset[0]
	}
	return current, false
}

func indexOf(arr []rune, r rune) int {
	for i, v := range arr {
		if v == r {
			return i
		}
	}
	return -1
}

func SaveState(filename string, state interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(state)
}

func LoadState(filename string, state interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	return decoder.Decode(state)
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)


func MonitorESC(escCh chan<- struct{}) {
	const VK_ESCAPE = 0x1B
	for {
		time.Sleep(50 * time.Millisecond)
		if isKeyPressed(VK_ESCAPE) {
			escCh <- struct{}{}
			return
		}
	}
}

func isKeyPressed(vkCode int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))
	return ret&0x8000 != 0
}

func SetupAutoRun() error {
	runKey := `Software\Microsoft\Windows\CurrentVersion\Run`

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %v", err)
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("cannot open registry key: %v", err)
	}
	defer k.Close()

	err = k.SetStringValue("BruteStick", exePath)
	if err != nil {
		return fmt.Errorf("cannot set registry value: %v", err)
	}

	return nil
}
