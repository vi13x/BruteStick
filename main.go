package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"bruteforce/autostart"
	"bruteforce/config"
	"bruteforce/core"
	"github.com/eiannone/keyboard"
)

var (
	stopFlag = false
)

func monitorEscapeKey() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Keyboard init error:", err)
		return
	}
	defer keyboard.Close()

	fmt.Println("Press ESC to stop")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		if key == keyboard.KeyEsc || char == 27 {
			stopFlag = true
			fmt.Println("\n[!] Stopped by user")
			return
		}

		runtime.Gosched()
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	go monitorEscapeKey()

	if autostart.IsRunningFromUSB() {
		autostart.SetupAutorun()
		exec.Command("cmd", "/C", "start", "/min", filepath.Base(os.Args[0])).Run()
		os.Exit(0)
	}

	fmt.Println("BruteStick Password Cracker")
	fmt.Println("--------------------------")

	state, err := core.LoadState()
	if err != nil {
		fmt.Printf("State load error: %v\n", err)
		return
	}

	var choice string
	if state.CurrentCombo != "" {
		fmt.Printf("Found saved state: length %d, combo '%s'\n",
			state.CurrentLength, state.CurrentCombo)
		fmt.Print("Continue (c) or start new (n)? [c/n]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			choice = scanner.Text()
		}
	}

	buffer := make([]rune, config.MaxLength)
	for length := state.CurrentLength; length <= config.MaxLength; length++ {
		if stopFlag {
			core.SaveState(length, string(buffer[:length]))
			return
		}

		fmt.Printf("\n[+] Brute-forcing length %d...\n", length)

		startFrom := ""
		if strings.ToLower(choice) != "n" && length == state.CurrentLength {
			startFrom = state.CurrentCombo
		}

		if !core.BruteForce(buffer, 0, length, startFrom, &stopFlag) {
			break
		}

		if err := core.SaveState(length+1, ""); err != nil {
			fmt.Printf("Save error: %v\n", err)
		}
	}

	fmt.Println("\n[!] All combinations exhausted")
}
