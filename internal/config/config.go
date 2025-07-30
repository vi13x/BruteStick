package config

import (
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

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	flag.IntVar(&cfg.MaxPasswordLen, "max-len", 6, "Максимальная длина пароля")
	flag.StringVar(&cfg.Alphabet, "alphabet", "abcdefghijklmnopqrstuvwxyz0123456789", "Алфавит для перебора")
	flag.StringVar(&cfg.StateFilePath, "state", "state.enc", "Путь к зашифрованному файлу состояния")
	flag.StringVar(&cfg.StateEncryptionKey, "key", os.Getenv("BRUTESTICK_KEY"), "Ключ шифрования состояния (или через env BRUTESTICK_KEY)")
	flag.StringVar(&cfg.LogPath, "mylog", "brutestick.mylog", "Путь к файлу лога")
	flag.BoolVar(&cfg.EnableAutostart, "autostart", false, "Включить автозапуск")
	flag.StringVar(&cfg.ExecutablePath, "exe-path", "", "Путь к исполняемому файлу (нужно для автозапуска)")
	flag.Parse()

	if cfg.StateEncryptionKey == "" {
		return nil, ErrMissingEncryptionKey
	}

	return cfg, nil
}

var ErrMissingEncryptionKey = &ConfigError{"Missing encryption key"}

type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string {
	return e.msg
}
