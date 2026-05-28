package main

import (
	"samrs-backend/internal/config"
	"samrs-backend/internal/seed"

	"go.uber.org/zap"
)

func main() {
	logger, cleanup := config.NewLogger()
	defer cleanup()
	_ = zap.RedirectStdLog(logger)

	db := config.ConnectDatabase()
	if err := seed.Run(db); err != nil {
		logger.Sugar().Fatalf("Seeder gagal: %v", err)
	}
	logger.Info("Seeder selesai.")
}
