package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vi13x/BruteStick/internal/brute"
	"github.com/vi13x/BruteStick/internal/config"
	"github.com/vi13x/BruteStick/internal/logger"
	"github.com/vi13x/BruteStick/internal/registry"
	"github.com/vi13x/BruteStick/internal/state"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg, err := logger.NewLogger(cfg.LogPath)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer func() {
		if err := logg.Close(); err != nil {
			log.Printf("Error closing logger: %v", err)
		}
	}()

	if cfg.EnableAutostart && cfg.ExecutablePath != "" {
		if err := registry.SetupAutostart(cfg.ExecutablePath); err != nil {
			logg.Error("Failed to setup autostart: %v", err)
		} else {
			logg.Info("Autostart configured")
		}
	}

	key := []byte(cfg.StateEncryptionKey)
	bruteState, err := state.LoadState(cfg.StateFilePath, key)
	if err != nil {
		logg.Warn("No previous state loaded, starting fresh")
		bruteState = state.NewBruteState()
	}

	// Отслеживание сигналов прерывания для корректной остановки
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		logg.Info("Interrupt received, stopping...")
		bruteState.Stop()
	}()

	// Запуск перебора
	brute.Run(bruteState, cfg, logg)

	// Сохранение состояния
	if err := state.SaveState(bruteState, cfg.StateFilePath, key); err != nil {
		logg.Error("Failed to save state: %v", err)
	}
}
