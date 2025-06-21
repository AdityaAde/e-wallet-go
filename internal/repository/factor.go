package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type factorRepository struct {
	db *goqu.Database
}

func NewFactor(con *sql.DB) domain.FactorRepository {
	return &factorRepository{
		db: goqu.New("default", con),
	}
}

// FindByUser implements domain.FactorRepository.
func (f factorRepository) FindByUser(ctx context.Context, id int64) (factor domain.Factor, err error) {
	dataset := f.db.From("factors").Where(goqu.Ex{
		"user_id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &factor)
	return
}
