package model

import (
	invoices_service "github.com/fidesy-pay/email-service/pkg/invoices-service"
	"github.com/google/uuid"
	"time"
)

type InvoiceMessage struct {
	ID             uuid.UUID                      `json:"id"`
	ClientID       uuid.UUID                      `json:"client_id"`
	UsdCentsAmount int64                          `json:"usd_cents_amount"`
	TokenAmount    *float64                       `json:"token_amount"`
	Chain          string                         `json:"chain"`
	Token          string                         `json:"token"`
	Status         invoices_service.InvoiceStatus `json:"status"`
	Address        string                         `json:"address"`
	CreatedAt      time.Time                      `json:"created_at"`
	PayerClientID  *string                        `json:"payer_client_id"`
	GasLimit       *int                           `json:"gas_limit"`
}
