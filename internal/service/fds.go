package service

import (
	"context"
	"log"
	"time"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
)

type fdsService struct {
	ipCheckerService   domain.IpCheckerService
	loginLogRepository domain.LoginLogRepository
}

func NewFds(ipCheckerService domain.IpCheckerService,
	loginLogRepository domain.LoginLogRepository) domain.FdsService {
	return &fdsService{
		ipCheckerService:   ipCheckerService,
		loginLogRepository: loginLogRepository,
	}
}

// IsAuthorized implements domain.FdsService.
func (f fdsService) IsAuthorized(ctx context.Context, ip string, userId int64) bool {
	locationCheck, err := f.ipCheckerService.Query(ctx, ip)

	if err != nil || locationCheck == (dto.IpChecker{}) {
		return false
	}

	newAccess := domain.LoginLog{
		UserID:       userId,
		IsAuthorized: false,
		IpAddress:    ip,
		Timezone:     locationCheck.Timezone,
		Long:         locationCheck.Long,
		Lat:          locationCheck.Lat,
		AccessTime:   time.Now(),
	}

	lastLogin, err := f.loginLogRepository.FindLastAuthorized(ctx, userId)

	if err != nil {
		_ = f.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	if lastLogin == (domain.LoginLog{}) {
		newAccess.IsAuthorized = true
		err := f.loginLogRepository.Save(ctx, &newAccess)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		return true
	}

	distanceHour := newAccess.AccessTime.Sub(lastLogin.AccessTime)
	distanceChange := util.GetDistance(lastLogin.Lat, lastLogin.Long, newAccess.Lat, newAccess.Long)

	log.Printf("distanceHour: %f", distanceHour.Hours())
	log.Printf("distanceChange: %f", distanceChange)

	if (distanceChange / distanceHour.Hours()) > 400 {
		_ = f.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	newAccess.IsAuthorized = true
	_ = f.loginLogRepository.Save(ctx, &newAccess)
	return true
}
