// cmd/main.go
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Egor213/notifyBot/internal/config"
	"github.com/Egor213/notifyBot/internal/handler"
	"github.com/Egor213/notifyBot/internal/repository"
	"github.com/Egor213/notifyBot/internal/service"
	"github.com/Egor213/notifyBot/pkg/bot"
	errorsUtils "github.com/Egor213/notifyBot/pkg/errors"
	"github.com/Egor213/notifyBot/pkg/logger"
	"github.com/Egor213/notifyBot/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

func Run() {
	// Config

	cfg, err := config.New()
	if err != nil {
		log.Fatal(errorsUtils.WrapPathErr(err))
	}

	// Logger
	logger.SetupLogger(cfg.Log.Level)
	log.Info("Logger has been set up")

	// Migrations
	Migrate(cfg.PG.URL)

	// DB connecting
	log.Info("Connecting to DB")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(errorsUtils.WrapPathErr(err))
	}
	defer pg.Close()
	log.Info("Connected to DB")

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не задан в .env или переменных окружения")
	}

	// Repo
	repo := repository.NewRepositoriesPG(pg)

	// Services
	services := service.NewServices(&service.ServiceDep{
		Repos: repo,
	})

	// Tg Bot
	ctx := context.Background()
	handler := handler.ConfigureHandler(ctx, services)
	telegramBot := bot.NewBot(token, handler, 10, true)
	telegramBot.Start(60)

	// Waiting signal
	log.Info("Configuring graceful shutdown")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	}

}
