package clientsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	"github.com/fidesy-pay/email-service/internal/pkg/model"
	"github.com/fidesy-pay/email-service/internal/pkg/templates"
)

type (
	Consumer struct {
		emailClient EmailClient
	}

	EmailClient interface {
		SendMessage(ctx context.Context, params dto.SendMessageParams) error
	}
)

func NewConsumer(
	emailClient EmailClient,
) *Consumer {
	return &Consumer{
		emailClient: emailClient,
	}
}

func (c *Consumer) Consume(ctx context.Context, msg []byte) error {
	client := new(model.ClientMessage)
	err := json.Unmarshal(msg, &client)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if !client.CreatedAt.After(time.Now().Add(-2 * time.Second)) {
		return nil
	}

	htmlTemplate, err := templates.BuildWelcome(client.Username)
	if err != nil {
		return fmt.Errorf("templates.BuildInvoiceComplete: %w", err)
	}

	err = c.emailClient.SendMessage(ctx, dto.SendMessageParams{
		SendTo:   client.Email,
		Subject:  "Thank you for registration!",
		HTMLBody: htmlTemplate,
	})
	if err != nil {
		return fmt.Errorf("emailClient.SendMessage: %w", err)
	}

	return nil
}
