package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type notificationRepository struct {
	db *goqu.Database
}

func NewNotification(con *sql.DB) domain.NotificationRepository {
	return &notificationRepository{
		db: goqu.New("default", con),
	}
}

// FindByUser implements domain.NotificationRepository.
func (n notificationRepository) FindByUser(ctx context.Context, user int64) (notifications []domain.Notification, err error) {
	dataset := n.db.From("notifications").Where(goqu.Ex{
		"user_id": user,
	}).Order(goqu.I("created_at").Desc()).Limit(15)

	err = dataset.ScanStructsContext(ctx, &notifications)
	return

}

// Insert implements domain.NotificationRepository.
func (n notificationRepository) Insert(ctx context.Context, notification *domain.Notification) error {
	executor := n.db.From("notifications").Insert().Rows(goqu.Record{
		"user_id":    notification.UserID,
		"title":      notification.Title,
		"body":       notification.Body,
		"status":     notification.Status,
		"is_read":    notification.IsRead,
		"created_at": notification.CreatedAt,
	}).Returning("id").Executor()

	_, err := executor.ScanStructContext(ctx, notification)
	return err
}

// Update implements domain.NotificationRepository.
func (n notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	executor := n.db.Update("notifications").Where(goqu.Ex{
		"id": notification.ID,
	}).Set(goqu.Record{
		"is_read": notification.IsRead,
	}).Executor()

	_, err := executor.ExecContext(ctx)
	return err

}
