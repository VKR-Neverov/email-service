package emailservice

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/fidesy-pay/email-service/internal/config"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	"github.com/fidesy-pay/email-service/internal/pkg/model"
	"github.com/fidesy-pay/email-service/internal/pkg/templates"
	"github.com/fidesy/sdk/common/postgres"
	"math/big"
)

const (
	codeLength = 6
	charset    = "0123456789"
)

type (
	Service struct {
		storage     Storage
		emailClient EmailClient
	}

	Storage interface {
		ListEmailCodes(ctx context.Context, params dto.ListEmailCodesParams) ([]*model.EmailCode, error)
		CreateEmailCode(ctx context.Context, emailCode *model.EmailCode) (*model.EmailCode, error)
		SetEmailCodeConfirmedAt(ctx context.Context, id string) (*model.EmailCode, error)
	}

	EmailClient interface {
		SendMessage(ctx context.Context, params dto.SendMessageParams) error
	}
)

func New(
	storage Storage,
	emailClient EmailClient,
) *Service {
	return &Service{
		storage:     storage,
		emailClient: emailClient,
	}
}

func (s *Service) SendCode(ctx context.Context, params dto.SendCodeParams) (*dto.SendCodeResult, error) {

	code := generateCode(codeLength)

	htmlTemplate, err := templates.BuildEmailCode(code)
	if err != nil {
		return nil, fmt.Errorf("templates.BuildEmailCode: %w", err)
	}

	err = s.emailClient.SendMessage(ctx, dto.SendMessageParams{
		SendTo:   params.Email,
		Subject:  "Verify your Fidesy.tech account",
		HTMLBody: htmlTemplate,
	})
	if err != nil {
		return nil, fmt.Errorf("emailClient.SendMessage: %w", err)
	}

	emailCode, err := s.storage.CreateEmailCode(ctx, &model.EmailCode{
		Code:   code,
		SentTo: params.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("storage.CreateEmailCode: %w", err)
	}

	return &dto.SendCodeResult{
		ID: emailCode.ID.String(),
	}, nil
}

func (s *Service) ConfirmCode(ctx context.Context, params dto.ConfirmCodeParams) (*dto.ConfirmCodeResult, error) {
	emailCodes, err := s.storage.ListEmailCodes(ctx, dto.ListEmailCodesParams{
		Filter: dto.ListEmailCodesFilter{
			IDIn: []string{params.ID},
		},
		Pagination: postgres.NewPagination(1, 1),
	})
	if err != nil {
		return nil, fmt.Errorf("storage.ListEmailCodes: %w", err)
	}

	if len(emailCodes) == 0 {
		return nil, ErrCodeNotFound
	}

	emailCode := emailCodes[0]

	// for staging testing
	whitelistCode := config.Get(config.WhiteListCode).(string)
	if params.Code != emailCode.Code && params.Code != whitelistCode {
		return nil, ErrCodeIsInvalid
	}

	_, err = s.storage.SetEmailCodeConfirmedAt(ctx, emailCode.ID.String())
	if err != nil {
		return nil, fmt.Errorf("storage.SetEmailCodeConfirmedAt: %w", err)
	}

	return &dto.ConfirmCodeResult{}, nil
}

func generateCode(length int) string {
	charsetLen := big.NewInt(int64(len(charset)))

	code := make([]byte, length)
	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(err)
		}
		code[i] = charset[randomIndex.Int64()]
	}

	return string(code)
}
