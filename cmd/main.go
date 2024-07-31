package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/encountea/message-service/config"
	"github.com/encountea/message-service/internal/handler"
	"github.com/encountea/message-service/internal/kafka"
	"github.com/encountea/message-service/internal/repository"
	"github.com/encountea/message-service/internal/service"
	"github.com/encountea/message-service/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load config: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to connect to db: %s", err.Error())
	}

	if err := repository.MigrateDB(db, "migrations"); err != nil {
		logrus.Fatalf("Error running migrations: %v", err)
	}

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logrus.Fatalf("Error creating Kafka producer: %v", err)
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo, kafkaProducer)
	handler := handler.NewHandler(service)

	app := new(server.Server)

	mux := handler.InitRoutes()

	go func() {
		if err := app.Run(cfg.Server.Port, mux); err != nil {
			logrus.Fatalf("Error to connect to server: %s", err.Error())
		}
	}()
	logrus.Printf("Server is running at port: %v", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server Shutting Down")
	if err := app.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on app shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
