package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"bruteforce/autostart"
	"bruteforce/config"
	"bruteforce/core"
	"github.com/eiannone/keyboard"
)

var stopFlag = false

func checkStopCombination() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Ошибка инициализации клавиатуры:", err)
		return
	}
	defer keyboard.Close()

	fmt.Println("Для остановки нажмите Ctrl+Пробел")

	var (
		ctrlPressed  bool
		spacePressed bool
		lastKey      keyboard.Key
	)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		// Определяем нажатие/отпускание Ctrl
		if key == keyboard.KeyCtrlL || key == keyboard.KeyCtrlR {
			if key == lastKey {
				// Клавиша удерживается - ничего не делаем
				continue
			}
			ctrlPressed = true
			lastKey = key
			continue
		}

		// Определяем нажатие/отпускание Space
		if key == keyboard.KeySpace {
			if key == lastKey {
				continue
			}
			spacePressed = true
			lastKey = key

			// Проверяем комбинацию
			if ctrlPressed && spacePressed {
				stopFlag = true
				fmt.Println("\n[!] Программа остановлена пользователем")
				os.Exit(0)
			}
			continue
		}

		// Если нажата другая клавиша - сбрасываем
		if key != lastKey {
			switch lastKey {
			case keyboard.KeyCtrlL, keyboard.KeyCtrlR:
				ctrlPressed = false
			case keyboard.KeySpace:
				spacePressed = false
			}
			lastKey = 0
		}

		time.Sleep(50 * time.Millisecond)
	}
}
func main() {
	go checkStopCombination()

	// Проверка запуска с USB и настройка автозапуска
	if autostart.IsRunningFromUSB() {
		autostart.SetupAutorun()
		exec.Command("cmd", "/C", "start", "/min", filepath.Base(os.Args[0])).Run()
		os.Exit(0)
	}

	fmt.Println("Brute Force Password Cracker")
	fmt.Println("--------------------------")

	// Остальной код без изменений
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
		if stopFlag {
			break
		}

		fmt.Printf("\n[+] Перебор паролей длины %d...\n", length)

		startFrom := ""
		if strings.ToLower(choice) != "n" && length == state.CurrentLength {
			startFrom = state.CurrentCombo
		}

		if !core.BruteForce(buffer, 0, length, startFrom, &stopFlag) {
			break
		}

		state.CurrentCombo = ""
		if err := core.SaveState(length+1, ""); err != nil {
			fmt.Printf("Ошибка сохранения состояния: %v\n", err)
		}
	}

	if !stopFlag {
		fmt.Println("\n[!] Все комбинации перебраны. Пароль не найден.")
	}
}
