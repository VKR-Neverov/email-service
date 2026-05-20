package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailCode struct {
	ID          uuid.UUID  `db:"id"`
	Code        string     `db:"code"`
	SentTo      string     `db:"sent_to"`
	SentAt      time.Time  `db:"sent_at"`
	ConfirmedAt *time.Time `db:"confirmed_at"`
}

func (e *EmailCode) TableName() string {
	return "email_codes"
}
