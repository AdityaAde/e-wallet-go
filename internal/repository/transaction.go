package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type transactionRepository struct {
	db *goqu.Database
}

func NewTransaction(con *sql.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: goqu.New("default", con),
	}
}

// Insert implements domain.TransactionRepository.
func (t transactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) error {

	executor := t.db.From("transactions").Insert().Rows(goqu.Record{
		"account_id":           transaction.AccountID,
		"sof_number":           transaction.SofNumber,
		"dof_number":           transaction.DofNumber,
		"amount":               transaction.Amount,
		"transaction_type":     transaction.TransactionType,
		"transaction_datetime": transaction.TransactionDateTime,
	}).Returning("id").Executor()

	_, err := executor.ScanStructContext(ctx, transaction)
	return err
}
