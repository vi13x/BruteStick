package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	MaxPasswordLen     int
	Alphabet           string
	StateFilePath      string
	StateEncryptionKey string
	LogPath            string
	EnableAutostart    bool
	ExecutablePath     string
}

var ErrMissingEncryptionKey = errors.New("missing encryption key")

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	flag.IntVar(&cfg.MaxPasswordLen, "max-len", 6, "Максимальная длина пароля")
	flag.StringVar(&cfg.Alphabet, "alphabet", "abcdefghijklmnopqrstuvwxyz0123456789", "Алфавит для перебора")
	flag.StringVar(&cfg.StateFilePath, "state", "state.enc", "Путь к зашифрованному файлу состояния")
	flag.StringVar(&cfg.StateEncryptionKey, "key", "", "Ключ шифрования состояния (или через env BRUTESTICK_KEY)")
	flag.StringVar(&cfg.LogPath, "log", "brutestick.log", "Путь к файлу лога")
	flag.BoolVar(&cfg.EnableAutostart, "autostart", false, "Включить автозапуск")
	flag.StringVar(&cfg.ExecutablePath, "exe-path", "", "Путь к исполняемому файлу (нужно для автозапуска)")

	flag.Parse()

	if cfg.StateEncryptionKey == "" {
		cfg.StateEncryptionKey = os.Getenv("BRUTESTICK_KEY")
		if cfg.StateEncryptionKey == "" {
			return nil, ErrMissingEncryptionKey
		}
	}

	// Ключ должен быть длиной 16, 24 или 32 байта для AES
	keyLen := len(cfg.StateEncryptionKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, errors.New("encryption key length must be 16, 24, or 32 bytes")
	}

	return cfg, nil
}
