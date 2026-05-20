package app

import (
	"context"
	"errors"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	emailservice "github.com/fidesy-pay/email-service/internal/pkg/email-service"
	desc "github.com/fidesy-pay/email-service/pkg/email-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ConfirmCode(ctx context.Context, req *desc.ConfirmCodeRequest) (*desc.ConfirmCodeResponse, error) {
	if err := validateConfirmCodeRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	_, err := i.emailService.ConfirmCode(ctx, dto.ConfirmCodeParams{
		ID:   req.GetId(),
		Code: req.GetCode(),
	})
	if err != nil {
		if errors.Is(err, emailservice.ErrCodeNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, emailservice.ErrCodeIsInvalid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "emailService.SendCode: %v", err)
	}

	return &desc.ConfirmCodeResponse{}, nil
}

func validateConfirmCodeRequest(req *desc.ConfirmCodeRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required, is.UUIDv4),
		validation.Field(&req.Code, validation.Required),
	)
}
