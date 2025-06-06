package domain

import (
	"context"
	"time"

	"adityaad.id/belajar-auth/dto"
)

type Transaction struct {
	ID                  int64     `json:"id"`
	AccountID           int64     `json:"account_id"`
	SofNumber           string    `json:"sof_number"`
	DofNumber           string    `json:"dof_number"`
	TransactionType     string    `json:"transaction_type"`
	Amount              float64   `json:"amount"`
	TransactionDateTime time.Time `json:"transaction_datetime"`
}

type TransactionRepository interface {
	Insert(ctx context.Context, transaction *Transaction) error
}

type TransactionService interface {
	TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error)
	TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error
}
