package service

import (
	"context"
	"errors"
	"fmt"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"github.com/google/uuid"
)

type topupService struct {
	notificationService domain.NotificationService
	topupRepository     domain.TopupRepository
	midtransService     domain.MidtransService
	accountRepository   domain.AccountRepository
}

func NewTopupService(notificationService domain.NotificationService, topupRepository domain.TopupRepository, midtransService domain.MidtransService,
	accountRepository domain.AccountRepository) domain.TopupService {
	return &topupService{
		notificationService: notificationService,
		topupRepository:     topupRepository,
		midtransService:     midtransService,
		accountRepository:   accountRepository,
	}
}

// InitializeTopup implements domain.TopupService.
func (t topupService) InitializeTopup(ctx context.Context, req dto.TopupReq) (dto.TopupRes, error) {
	topup := domain.Topup{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Status: 0,
		Amount: req.Amount,
	}

	err := t.midtransService.GenerateSnapURL(ctx, &topup)

	if err != nil {
		return dto.TopupRes{}, err
	}

	err = t.topupRepository.Insert(ctx, &topup)
	if err != nil {
		return dto.TopupRes{}, err
	}

	return dto.TopupRes{
		SnapURL: topup.SnapURL,
	}, nil
}

// ConfirmedTopup implements domain.TopupService.
func (t topupService) ConfirmedTopup(ctx context.Context, id string) error {

	topup, err := t.topupRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if topup == (domain.Topup{}) {
		return errors.New("topup not found")
	}

	account, err := t.accountRepository.FindByUserId(ctx, topup.UserID)
	if err != nil {
		return err
	}

	if account == (domain.Account{}) {
		return errors.New("account not found")
	}

	account.Balance += topup.Amount
	err = t.accountRepository.Update(ctx, &account)

	if err != nil {
		return err
	}

	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", topup.Amount),
	}

	_ = t.notificationService.Insert(ctx, account.UserId, "TOPUP_SUCCESS", data)

	return err
}
