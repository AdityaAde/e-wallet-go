package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
)

type transactionService struct {
	accountRepository      domain.AccountRepository
	transactionRepository  domain.TransactionRepository
	cacheRepository        domain.CacheRepository
	notificationRepository domain.NotificationRepository
	hub                    *dto.Hub
}

func NewTransaction(accountRepository domain.AccountRepository, transactionRepository domain.TransactionRepository, cacheRepository domain.CacheRepository, notificationRepository domain.NotificationRepository, hub *dto.Hub) domain.TransactionService {
	return &transactionService{
		accountRepository:      accountRepository,
		transactionRepository:  transactionRepository,
		cacheRepository:        cacheRepository,
		notificationRepository: notificationRepository,
		hub:                    hub,
	}
}

// TransferInquiry implements domain.TransactionService.
func (t transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := t.accountRepository.FindByUserId(ctx, user.ID)

	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)

	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound

	}

	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := util.GenerateRandomString(32)
	jsonData, _ := json.Marshal(req)

	_ = t.cacheRepository.Set(inquiryKey, jsonData)

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

// TransferExecute implements domain.TransactionService.
func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	val, err := t.cacheRepository.Get(req.InquiryKey)

	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	_ = json.Unmarshal(val, &reqInq)

	if reqInq == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := t.accountRepository.FindByUserId(ctx, user.ID)

	if err != nil {
		return err
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, reqInq.AccountNumber)

	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction{
		AccountID:           myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "D",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	creditTransaction := domain.Transaction{
		AccountID:           dofAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "C",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &creditTransaction)
	if err != nil {
		return err
	}

	myAccount.Balance -= reqInq.Amount
	err = t.accountRepository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += reqInq.Amount
	err = t.accountRepository.Update(ctx, &dofAccount)
	if err != nil {
		return err
	}

	go t.NotificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)

	return nil
}

func (t transactionService) NotificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	notificationSender := domain.Notification{
		UserID:    sofAccount.UserId,
		Title:     "Transfer Berhasil",
		Body:      fmt.Sprintf("Transfer sebesar %.2f berhasil", amount),
		Status:    0,
		IsRead:    1,
		CreatedAt: time.Now(),
	}

	notificationReceiver := domain.Notification{
		UserID:    dofAccount.UserId,
		Title:     "Dana Masuk",
		Body:      fmt.Sprintf("Dana masuk senilai %.2f", amount),
		Status:    0,
		IsRead:    1,
		CreatedAt: time.Now(),
	}

	_ = t.notificationRepository.Insert(context.Background(), &notificationSender)

	if channel, ok := t.hub.NotificationChannel[sofAccount.UserId]; ok {
		channel <- dto.NotificationData{
			ID:        notificationSender.ID,
			Title:     notificationSender.Title,
			Body:      notificationSender.Body,
			Status:    notificationSender.Status,
			IsRead:    notificationSender.IsRead,
			CreatedAt: notificationSender.CreatedAt,
		}
	}

	_ = t.notificationRepository.Insert(context.Background(), &notificationReceiver)

	if channel, ok := t.hub.NotificationChannel[dofAccount.UserId]; ok {
		channel <- dto.NotificationData{
			ID:        notificationReceiver.ID,
			Title:     notificationReceiver.Title,
			Body:      notificationReceiver.Body,
			Status:    notificationReceiver.Status,
			IsRead:    notificationReceiver.IsRead,
			CreatedAt: notificationReceiver.CreatedAt,
		}
	}

}
