package app

import (
	"context"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	desc "github.com/fidesy-pay/email-service/pkg/email-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SendCode(ctx context.Context, req *desc.SendCodeRequest) (*desc.SendCodeResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	result, err := i.emailService.SendCode(ctx, dto.SendCodeParams{
		Email: req.Email,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "emailService.SendCode: %v", err)
	}

	return &desc.SendCodeResponse{
		Id: result.ID,
	}, nil
}
