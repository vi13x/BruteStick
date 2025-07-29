package core

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"bruteforce/config"
	"github.com/go-vgo/robotgo"
)

func BruteForce(buffer []rune, pos, length int, startFrom string) bool {
	if pos == length {
		password := string(buffer[:length])
		fmt.Printf("Пробуем: %s\n", password)
		saveResult(password)
		sendKeys(password)
		return false
	}

	if pos < len(startFrom) {
		buffer[pos] = rune(startFrom[pos])
		if BruteForce(buffer, pos+1, length, startFrom) {
			return true
		}
		return false
	}

	for _, char := range config.Alphabet {
		buffer[pos] = char
		if BruteForce(buffer, pos+1, length, startFrom) {
			return true
		}
	}
	return false
}

func saveResult(password string) {
	file, err := os.OpenFile(config.ResultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "%s\n", password)
}

func sendKeys(keys string) {
	upperKeys := ""
	for _, r := range keys {
		upperKeys += string(unicode.ToUpper(r))
	}
	robotgo.TypeStr(upperKeys)
	robotgo.KeyTap("enter")
	time.Sleep(time.Millisecond * config.DelayBetween)
}
