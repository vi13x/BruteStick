package core

import (
	"brutestick/internal/config"
	"brutestick/internal/logger"
	"brutestick/internal/utils"
	"fmt"
	"sync"
)

type BruteState struct {
	CurrentPassword string
}

func Run(conf *config.Config, log *logger.Logger, escCh <-chan struct{}) error {
	log.Info("Core started with max password length %d", conf.MaxPasswordLength)

	state := &BruteState{}
	err := utils.LoadState(conf.SaveFile, state)
	if err != nil {
		log.Warn("No previous state found, starting new brute force")
		state.CurrentPassword = ""
	} else {
		log.Info("Loaded previous state: %s", state.CurrentPassword)
	}

	charSet := utils.DefaultCharSet()
	maxLen := conf.MaxPasswordLength

	var mu sync.Mutex
	stopped := false

	go func() {
		<-escCh
		mu.Lock()
		stopped = true
		mu.Unlock()
	}()

	password := state.CurrentPassword

	log.Info("Brute forcing started from: %s", password)

	for length := len(password); length <= maxLen; length++ {
		if length == 0 {
			length = 1
		}
		passwordRunes := []rune(password)
		passwordRunes = utils.PadPassword(passwordRunes, length, charSet[0])
		for {
			mu.Lock()
			if stopped {
				mu.Unlock()
				log.Info("Brute forcing stopped by user")
				return nil
			}
			mu.Unlock()

			fmt.Printf("Trying password: %s\r", string(passwordRunes))

			state.CurrentPassword = string(passwordRunes)
			err := utils.SaveState(conf.SaveFile, state)
			if err != nil {
				log.Warn("Failed to save state: %v", err)
			}

			next, ok := utils.NextPassword(passwordRunes, charSet)
			if !ok {
				break
			}
			passwordRunes = next
		}
		password = ""
	}

	log.Info("Brute forcing completed all combinations")
	return nil
}
