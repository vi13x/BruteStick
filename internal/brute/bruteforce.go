package brute

import (
	"strings"

	"github.com/vi13x/BruteStick/internal/config"
	"github.com/vi13x/BruteStick/internal/logger"
	"github.com/vi13x/BruteStick/internal/state"
)

// Run запускает цикл перебора паролей
func Run(s *state.BruteState, cfg *config.Config, logg *logger.Logger) {
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
			logg.Info("Trying password: %s", pwd)

			// TODO: Добавьте логику проверки пароля здесь (например, попытка аутентификации)

			pwd, done = nextPassword(pwd, alphabet)
			s.CurrentPassword = pwd
		}

		s.CurrentLength++
	}
}

// nextPassword генерирует следующий пароль. Возвращает следующий пароль и true, если перебор закончился
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
				// Если достигнут последний пароль данной длины, сигнализируем завершение текущей длины
				return string(bytes), true
			}
		} else {
			bytes[i] = alphabet[index]
			break
		}
	}
	return string(bytes), false
}
