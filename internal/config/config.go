package config

import (
	"flag"
)

type Config struct {
	MaxPasswordLength int
	SaveFile          string
}

func LoadConfig() *Config {
	maxLen := flag.Int("maxlen", 6, "Максимальная длина пароля для перебора")
	saveFile := flag.String("savefile", "state.dat", "Файл для сохранения прогресса")
	flag.Parse()

	return &Config{
		MaxPasswordLength: *maxLen,
		SaveFile:          *saveFile,
	}
}
