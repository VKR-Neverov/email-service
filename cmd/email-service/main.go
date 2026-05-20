package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fidesy-pay/email-service/internal/app"
	"github.com/fidesy-pay/email-service/internal/config"
	clientsconsumer "github.com/fidesy-pay/email-service/internal/pkg/consumers/clients-consumer"
	"github.com/fidesy-pay/email-service/internal/pkg/consumers/invoice-consumer"
	emailclient "github.com/fidesy-pay/email-service/internal/pkg/email-client"
	emailservice "github.com/fidesy-pay/email-service/internal/pkg/email-service"
	"github.com/fidesy-pay/email-service/internal/pkg/storage"
	clients_service "github.com/fidesy-pay/email-service/pkg/clients-service"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/kafka"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	server, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "domain-name-service:10000"),
		grpc.WithGraylog("graylog:5555"),
		grpc.WithTracer("http://jaeger:14268/api/traces"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	clientsClient, err := grpc.NewClient[clients_service.ClientsServiceClient](
		ctx,
		clients_service.NewClientsServiceClient,
		"rpc:///clients-service",
	)
	if err != nil {
		logger.Fatalf("NewClientsClient: %v", err)
	}

	pool, err := postgres.Connect(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		logger.Fatalf("postgres.Connect: %v", err)
	}

	storage := storage.New(pool)

	emailClient := emailclient.New()

	err = kafka.RegisterConsumer(
		ctx,
		invoiceconsumer.NewInvoiceConsumer(emailClient, clientsClient),
		config.Get(config.KafkaBrokers).([]string),
		"invoices-json",
	)
	if err != nil {
		log.Fatalf("kafka.RegisterConsumer: %v", err)
	}

	err = kafka.RegisterConsumer(
		ctx,
		clientsconsumer.NewConsumer(emailClient),
		config.Get(config.KafkaBrokers).([]string),
		"clients-json",
	)
	if err != nil {
		log.Fatalf("kafka.RegisterConsumer: %v", err)
	}

	emailService := emailservice.New(storage, emailClient)

	impl := app.New(emailService)

	err = server.Run(ctx, impl)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
