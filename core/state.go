package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"bruteforce/config"
)

type ProgressState struct {
	CurrentLength int
	CurrentCombo  string
}

func SaveState(length int, combo string) error {
	file, err := os.Create(config.StateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n%s", length, combo)
	return err
}

func LoadState() (ProgressState, error) {
	state := ProgressState{CurrentLength: config.MinLength}

	if _, err := os.Stat(config.StateFile); os.IsNotExist(err) {
		return state, nil
	}

	file, err := os.Open(config.StateFile)
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
