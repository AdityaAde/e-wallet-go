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
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
	notificationService   domain.NotificationService
}

func NewTransaction(accountRepository domain.AccountRepository, transactionRepository domain.TransactionRepository, cacheRepository domain.CacheRepository, notificationService domain.NotificationService) domain.TransactionService {
	return &transactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
		notificationService:   notificationService,
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
	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", amount),
	}

	_ = t.notificationService.Insert(context.Background(), sofAccount.UserId, "TRANSFER", data)
	_ = t.notificationService.Insert(context.Background(), dofAccount.UserId, "TRANSFER_DES", data)

}
