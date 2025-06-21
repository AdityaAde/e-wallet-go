package domain

import (
	"context"

	"adityaad.id/belajar-auth/dto"
)

type Topup struct {
	ID      string  `db:"id"`
	UserID  int64   `db:"user_id"`
	Status  int8    `db:"status"`
	Amount  float64 `db:"amount"`
	SnapURL string  `db:"snap_url"`
}

type TopupRepository interface {
	FindById(ctx context.Context, id string) (Topup, error)
	Insert(ctx context.Context, topup *Topup) error
	Update(ctx context.Context, topup *Topup) error
}

type TopupService interface {
	ConfirmedTopup(ctx context.Context, id string) error
	InitializeTopup(ctx context.Context, req dto.TopupReq) (dto.TopupRes, error)
}
