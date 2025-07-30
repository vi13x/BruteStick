package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"bruteforce/internal/brute"
	"bruteforce/internal/config"
	"bruteforce/internal/registry"
	"bruteforce/internal/state"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логирования
	logger, err := log.NewLogger(cfg.LogPath)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Close()

	// Запуск автозапуска в реестре (при условии флага --autostart)
	if cfg.EnableAutostart {
		if err := registry.SetupAutostart(cfg.ExecutablePath); err != nil {
			logger.Error("Failed to setup autostart: %v", err)
		} else {
			logger.Info("Autostart configured")
		}
	}

	// Загрузка состояния с шифрованием
	key := []byte(cfg.StateEncryptionKey)
	bruteState, err := state.LoadState(cfg.StateFilePath, key)
	if err != nil {
		logger.Warn("No previous state loaded, starting fresh")
		bruteState = state.NewBruteState()
	}

	// Запуск перебора с прерыванием по Ctrl+C
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		logger.Info("Interrupt received, stopping...")
		bruteState.Stop()
	}()

	// Основной цикл перебора
	brute.Run(bruteState, cfg, logger)

	// Сохранение состояния перед выходом
	if err := state.SaveState(bruteState, cfg.StateFilePath, key); err != nil {
		logger.Error("Failed to save state: %v", err)
	}
}
