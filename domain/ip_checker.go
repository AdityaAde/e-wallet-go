package domain

import (
	"context"

	"adityaad.id/belajar-auth/dto"
)

type IpCheckerService interface {
	Query(ctx context.Context, ip string) (dto.IpChecker, error)
}
