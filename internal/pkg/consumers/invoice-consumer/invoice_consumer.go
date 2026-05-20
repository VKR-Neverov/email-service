package invoiceconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	"github.com/fidesy-pay/email-service/internal/pkg/model"
	"github.com/fidesy-pay/email-service/internal/pkg/templates"
	clients_service "github.com/fidesy-pay/email-service/pkg/clients-service"
	invoices_service "github.com/fidesy-pay/email-service/pkg/invoices-service"
	"github.com/fidesy/sdk/common/logger"
	"google.golang.org/grpc"
)

type (
	InvoiceConsumer struct {
		emailClient   EmailClient
		clientsClient ClientsClient
	}

	EmailClient interface {
		SendMessage(ctx context.Context, params dto.SendMessageParams) error
	}

	ClientsClient interface {
		ListClients(ctx context.Context, in *clients_service.ListClientsRequest, opts ...grpc.CallOption) (*clients_service.ListClientsResponse, error)
	}
)

func NewInvoiceConsumer(
	emailClient EmailClient,
	clientsClient ClientsClient,
) *InvoiceConsumer {
	return &InvoiceConsumer{
		emailClient:   emailClient,
		clientsClient: clientsClient,
	}
}

func (c *InvoiceConsumer) Consume(ctx context.Context, msg []byte) error {
	invoice := new(model.InvoiceMessage)
	err := json.Unmarshal(msg, &invoice)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if invoice.Status != invoices_service.InvoiceStatus_SUCCESS {
		return nil
	}

	clients, err := c.clientsClient.ListClients(ctx, &clients_service.ListClientsRequest{
		Filter: &clients_service.ListClientsRequest_Filter{
			IdIn: []string{invoice.ClientID.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("clientsClient.ListClients: %w", err)
	}

	if len(clients.Clients) == 0 {
		logger.Info(fmt.Sprintf("client not found by id=%s", invoice.ClientID.String()))
		return nil
	}

	client := clients.Clients[0]

	if !client.IsInvoiceNotificationEnabled {
		return nil
	}

	htmlTemplate, err := templates.BuildInvoiceComplete(
		invoice.ID.String(),
		float64(invoice.UsdCentsAmount)/100,
	)
	if err != nil {
		return fmt.Errorf("templates.BuildInvoiceComplete: %w", err)
	}

	err = c.emailClient.SendMessage(ctx, dto.SendMessageParams{
		SendTo:   client.Email,
		Subject:  "Invoice has been paid",
		HTMLBody: htmlTemplate,
	})
	if err != nil {
		return fmt.Errorf("emailClient.SendMessage: %w", err)
	}

	return nil
}
