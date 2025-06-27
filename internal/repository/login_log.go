package repository

import (
	"context"
	"database/sql"

	"adityaad.id/belajar-auth/domain"
	"github.com/doug-martin/goqu/v9"
)

type LoginLogRepository struct {
	db *goqu.Database
}

func NewLoginLog(con *sql.DB) domain.LoginLogRepository {
	return &LoginLogRepository{
		db: goqu.New("default", con),
	}
}

// FindLastAuthorized implements domain.LoginLogRepository.
func (l LoginLogRepository) FindLastAuthorized(ctx context.Context, userId int64) (loginlog domain.LoginLog, err error) {
	dataset := l.db.From("login_log").Where(goqu.Ex{
		"user_id":       userId,
		"is_authorized": true,
	}).Order(goqu.I("access_time").Desc()).Limit(1)

	if _, err = dataset.ScanStructContext(ctx, &loginlog); err != nil {
		return domain.LoginLog{}, err
	}

	return

}

// Save implements domain.LoginLogRepository.
func (l LoginLogRepository) Save(ctx context.Context, login *domain.LoginLog) error {
	exec := l.db.Insert("login_log").Rows(goqu.Record{
		"user_id":       login.UserID,
		"is_authorized": login.IsAuthorized,
		"ip_address":    login.IpAddress,
		"timezone":      login.Timezone,
		"lat":           login.Lat,
		"lon":           login.Long,
		"access_time":   login.AccessTime,
	}).Returning("id").Executor()

	_, err := exec.ScanStructContext(ctx, login)
	return err
}
