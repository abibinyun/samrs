package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"samrs-backend/internal/app"
	"samrs-backend/internal/config"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 2. SETUP LOGGER
	logger, cleanup := config.NewLogger()
	defer cleanup()
	_ = zap.RedirectStdLog(logger)

	// 1. Inisialisasi Database & AutoMigrate
	// GORM akan otomatis sinkronisasi struct domain ke table PostgreSQL
	log := logger.Sugar()

	log.Info("Connecting to database...")
	db := config.ConnectDatabase()
	log.Info("Database connected")
	
	// 3. APP (DI + Router)
	appInstance := app.New(db, logger)

	// 5. JALANKAN SERVER
	addr := resolveHTTPAddr()
	server := &http.Server{
		Addr:    addr,
		Handler: appInstance.Router,
	}

	go func() {
		log.Infof("Server berjalan di %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal menjalankan server: %v", err)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-stopCtx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Warnf("Gagal shutdown server: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func resolveHTTPAddr() string {
	port := strings.TrimSpace(os.Getenv("APP_PORT"))
	if port == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}
