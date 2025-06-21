package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type topupRepository struct {
	db *goqu.Database
}

func NewTopup(con *sql.DB) domain.TopupRepository {
	return &topupRepository{
		db: goqu.New("default", con),
	}
}

// FindById implements domain.TopupRepository.
func (r topupRepository) FindById(ctx context.Context, id string) (topup domain.Topup, err error) {
	dataset := r.db.From("topup").Where(goqu.Ex{
		"id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &topup)
	return
}

// Insert implements domain.TopupRepository.
func (r topupRepository) Insert(ctx context.Context, topup *domain.Topup) error {
	executor := r.db.Insert("topup").Rows(goqu.Record{
		"id":       topup.ID,
		"user_id":  topup.UserID,
		"status":   topup.Status,
		"amount":   topup.Amount,
		"snap_url": topup.SnapURL,
	}).Executor()

	_, err := executor.ScanStructContext(ctx, topup)
	return err
}

// Update implements domain.TopupRepository.
func (r topupRepository) Update(ctx context.Context, topup *domain.Topup) error {
	executor := r.db.Update("topup").Where(goqu.Ex{
		"id": topup.ID,
	}).Set(goqu.Record{
		"status":   topup.Status,
		"amount":   topup.Amount,
		"snap_url": topup.SnapURL,
	}).Executor()

	_, err := executor.ScanStructContext(ctx, topup)
	return err
}
