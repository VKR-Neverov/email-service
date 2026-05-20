package emailclient

import (
	"context"
	"crypto/tls"
	"github.com/fidesy-pay/email-service/internal/config"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	gomail "gopkg.in/mail.v2"
	"os"
)

type Client struct {
	client *gomail.Dialer
}

func New() *Client {
	d := gomail.NewDialer(
		"smtp.mail.ru",
		465,
		"mail@fidesy.tech",
		config.Get(config.EmailPassword).(string),
	)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		client: d,
	}
}

func (c *Client) SendMessage(_ context.Context, params dto.SendMessageParams) error {
	// Send code only on production
	env := os.Getenv("ENV")
	if env != "PRODUCTION" {
		return nil
	}

	m := gomail.NewMessage()

	m.SetHeader("From", "Fidesy Tech <mail@fidesy.tech>")
	m.SetHeader("To", params.SendTo)
	m.SetHeader("Subject", params.Subject)
	m.SetBody("text/html", params.HTMLBody)

	if err := c.client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
