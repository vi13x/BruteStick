package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"brutestick/internal/config"
	"brutestick/internal/core"
	"brutestick/internal/logger"
	"brutestick/internal/utils"
)

func main() {
	log := logger.NewLogger()
	defer log.Close()

	conf := config.LoadConfig()
	log.Info("Starting brutestick with config: %+v", conf)

	err := utils.SetupAutoRun()
	if err != nil {
		log.Warn("Failed to setup autorun: %v", err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	escCh := make(chan struct{})
	go utils.MonitorESC(escCh)

	doneCh := make(chan error)
	go func() {
		doneCh <- core.Run(conf, log, escCh)
	}()

	select {
	case <-stopCh:
		log.Info("Termination signal received, exiting...")
		fmt.Println("\nExit signal received")
	case <-escCh:
		log.Info("ESC pressed, stopping brute force")
		fmt.Println("\nESC pressed, stopping brute force")
	case err := <-doneCh:
		if err != nil {
			log.Error("Error in brute force: %v", err)
			fmt.Println("Error:", err)
		} else {
			log.Info("Brute force finished successfully")
			fmt.Println("Brute force finished successfully")
		}
	}

	time.Sleep(500 * time.Millisecond)
}
