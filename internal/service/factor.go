package service

import (
	"context"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"golang.org/x/crypto/bcrypt"
)

type factorService struct {
	factorRepository domain.FactorRepository
}

func NewFactor(factorRepository domain.FactorRepository) domain.FactorService {
	return &factorService{
		factorRepository: factorRepository,
	}
}

// ValidatePIN implements domain.FactorService.
func (f factorService) ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error {
	factor, err := f.factorRepository.FindByUser(ctx, req.UserID)
	if err != nil {
		return err
	}

	if factor == (domain.Factor{}) {
		return domain.ErrPinInvalid
	}

	err = bcrypt.CompareHashAndPassword([]byte(factor.Pin), []byte(req.PIN))
	if err != nil {
		return domain.ErrPinInvalid
	}

	return nil
}
