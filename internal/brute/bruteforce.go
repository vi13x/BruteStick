package brute

import (
	"bruteforce/internal/state"
	"log"
	"strings"
)

func Run(s *state.BruteState, cfg *Config, logger *log.Logger) {
	alphabet := cfg.Alphabet
	maxLen := cfg.MaxPasswordLen

	if s.Stopped {
		return
	}

	for length := s.CurrentLength; length <= maxLen; length++ {
		pwd := s.CurrentPassword
		if len(pwd) != length {
			pwd = strings.Repeat(string(alphabet[0]), length)
			s.CurrentPassword = pwd
		}

		done := false
		for !done && !s.Stopped {
			logger.Info("Trying password: %s", pwd)
			// Здесь добавить вызов попытки аутентификации или проверку пароля

			pwd, done = nextPassword(pwd, alphabet)
			s.CurrentPassword = pwd
		}

		s.CurrentLength++
	}
}

// nextPassword генерирует следующий пароль. Возвращает новый пароль и true, если перебор закончился.
func nextPassword(current, alphabet string) (string, bool) {
	bytes := []byte(current)
	base := len(alphabet)

	for i := len(bytes) - 1; i >= 0; i-- {
		index := strings.IndexRune(alphabet, rune(bytes[i]))
		if index == -1 {
			index = 0
		}
		index++
		if index == base {
			bytes[i] = alphabet[0]
			if i == 0 {
				return string(bytes), true
			}
		} else {
			bytes[i] = alphabet[index]
			break
		}
	}

	return string(bytes), false
}
