package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

// Insert implements domain.UserRepository.

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", con),
	}
}

func (u userRepository) Insert(ctx context.Context, user *domain.User) error {
	executor := u.db.From("users").Insert().Rows(goqu.Record{
		"full_name": user.FullName,
		"phone":     user.Phone,
		"email":     user.Email,
		"username":  user.Username,
		"password":  user.Password,
	}).Returning("id").Executor()

	_, err := executor.ScanStructContext(ctx, user)
	return err
}

// Update implements domain.UserRepository.
func (u userRepository) Update(ctx context.Context, user *domain.User) error {
	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	executor := u.db.Update("users").Set(goqu.Record{
		"full_name":         user.FullName,
		"phone":             user.Phone,
		"email":             user.Email,
		"username":          user.Username,
		"password":          user.Password,
		"email_verified_at": user.EmailVerifiedAtDB,
	}).Where(goqu.Ex{
		"id": user.ID,
	}).Executor()

	_, err := executor.ExecContext(ctx)
	return err
}

// FindById implements domain.UserRepository.
func (u userRepository) FindById(ctx context.Context, id int64) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &user)
	return
}

// FindByUsername implements domain.UserRepository.
func (u userRepository) FindByUsername(ctx context.Context, username string) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"username": username,
	})

	_, err = dataset.ScanStructContext(ctx, &user)
	return
}
