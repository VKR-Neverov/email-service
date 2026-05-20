package app

import (
	"context"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	desc "github.com/fidesy-pay/email-service/pkg/email-service"
	"google.golang.org/grpc"
)

type (
	Implementation struct {
		desc.UnimplementedEmailServiceServer

		emailService EmailService
	}

	EmailService interface {
		SendCode(ctx context.Context, params dto.SendCodeParams) (*dto.SendCodeResult, error)
		ConfirmCode(ctx context.Context, params dto.ConfirmCodeParams) (*dto.ConfirmCodeResult, error)
	}
)

func New(
	emailService EmailService,
) *Implementation {
	return &Implementation{
		emailService: emailService,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.EmailService_ServiceDesc
}
