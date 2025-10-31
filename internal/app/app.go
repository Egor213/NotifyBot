// cmd/main.go
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	kafkabroker "github.com/Egor213/notifyBot/internal/broker/kafka"
	telegramworker "github.com/Egor213/notifyBot/internal/broker/worker/telegram"
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
	cfg, err := config.New()
	if err != nil {
		log.Fatal(errorsUtils.WrapPathErr(err))
	}

	logger.SetupLogger(cfg.Log.Level)
	log.Info("Logger has been set up")

	Migrate(cfg.PG.URL)

	log.Info("Connecting to DB")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(errorsUtils.WrapPathErr(err))
	}
	defer pg.Close()
	log.Info("Connected to DB")

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не задан")
	}

	repo := repository.NewRepositoriesPG(pg)

	sendMailKey := os.Getenv("MAIL_KEY")
	if sendMailKey == "" {
		log.Fatal("MAIL_KEY не задан")
	}

	sendMailDep := service.SendMailDep{
		SendMail:  os.Getenv("SENDER_MAIL"),
		Port:      587,
		Protocol:  "smtp.gmail.com",
		SecretKey: sendMailKey,
	}
	services := service.NewServices(&service.ServiceDep{
		Repos:       repo,
		SendMailDep: sendMailDep,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := handler.ConfigureHandler(ctx, services)
	telegramBot := bot.NewBot(token, handler, 10, true)

	br := kafkabroker.NewConsumer(kafkabroker.ConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notifications",
		GroupID: "notify-bot",
	})
	kafkaWorker := telegramworker.NewNotifyWorker(telegramBot, services.NotifySettings)
	br.RegisterWorker(kafkaWorker)

	go func() {
		if err := br.Run(ctx); err != nil {
			log.Errorf("Kafka consumer stopped with error: %v", err)
		}
	}()

	go telegramBot.Start(60)

	// Graceful shutdown
	log.Info("Configuring graceful shutdown")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Info("Received shutdown signal")

	cancel()
	br.Close()
	log.Info("Shutdown complete")
}
