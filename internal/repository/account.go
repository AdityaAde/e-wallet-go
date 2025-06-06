package repository

import (
	"context"
	"database/sql"
	"log"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type accountRepository struct {
	db *goqu.Database
}

func NewAccount(con *sql.DB) domain.AccountRepository {
	return &accountRepository{
		db: goqu.New("default", con),
	}
}

// FindByUserId implements domain.AccountRepository.
func (a accountRepository) FindByUserId(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := a.db.From("accounts").Where(goqu.Ex{
		"user_id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

// FindByAccpountNumber implements domain.AccountRepository.
func (a accountRepository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := a.db.From("accounts").Where(goqu.Ex{
		"account_number": accNumber,
	})

	log.Println("find by account number")

	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

// Update implements domain.AccountRepository.
func (a accountRepository) Update(ctx context.Context, account *domain.Account) error {

	executor := a.db.Update("accounts").Where(goqu.Ex{
		"id": account.ID,
	}).Set(goqu.Record{
		"balance": account.Balance,
	}).Executor()

	_, err := executor.ExecContext(ctx)
	return err
}
