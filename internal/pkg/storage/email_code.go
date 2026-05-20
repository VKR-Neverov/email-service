package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy-pay/email-service/internal/pkg/dto"
	"github.com/fidesy-pay/email-service/internal/pkg/model"
	"github.com/fidesy/sdk/common/postgres"
	"time"
)

func (s *Storage) ListEmailCodes(ctx context.Context, params dto.ListEmailCodesParams) ([]*model.EmailCode, error) {
	query := postgres.Builder().
		Select(emailCodeFields).
		From(emailCodesTable)

	if len(params.Filter.IDIn) > 0 {
		query = query.Where(sq.Eq{
			"id": params.Filter.IDIn,
		})
	}

	query = query.
		Limit(params.Pagination.Limit()).
		Offset(params.Pagination.Offset())

	return postgres.Select[model.EmailCode](ctx, s.pool, query)
}

func (s *Storage) CreateEmailCode(ctx context.Context, emailCode *model.EmailCode) (*model.EmailCode, error) {
	query := postgres.Builder().
		Insert(emailCodesTable).
		SetMap(map[string]interface{}{
			"code":    emailCode.Code,
			"sent_to": emailCode.SentTo,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", emailCodeFields))

	return postgres.Exec[model.EmailCode](ctx, s.pool, query)
}

func (s *Storage) SetEmailCodeConfirmedAt(ctx context.Context, id string) (*model.EmailCode, error) {
	query := postgres.Builder().
		Update(emailCodesTable).
		SetMap(map[string]interface{}{
			"confirmed_at": time.Now(),
		}).
		Where(sq.Eq{
			"id": id,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", emailCodeFields))

	return postgres.Exec[model.EmailCode](ctx, s.pool, query)
}
