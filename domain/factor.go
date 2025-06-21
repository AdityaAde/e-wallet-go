package domain

import (
	"context"

	"adityaad.id/belajar-auth/dto"
)

type Factor struct {
	ID   int64  `db:"id"`
	Name string `db:"user_id"`
	Pin  string `db:"pin"`
}

type FactorRepository interface {
	FindByUser(ctx context.Context, user int64) (Factor, error)
}

type FactorService interface {
	ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error
}
