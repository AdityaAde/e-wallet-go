package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository  domain.UserRepository
	cacheRepository domain.CacheRepository
	emailService    domain.EmailService
}

// Register implements domain.UserService.

func NewUser(userRepository domain.UserRepository, cacheRepository domain.CacheRepository, emailService domain.EmailService) domain.UserService {
	return &userService{
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
		emailService:    emailService,
	}
}

// Authenticate implements domain.UserService.
func (u userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}

	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	if !user.EmailVerifiedAtDB.Valid {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(16)

	userJson, _ := json.Marshal(user)
	_ = u.cacheRepository.Set("user"+token, userJson)

	return dto.AuthRes{
		Token: token,
	}, nil
}

// ValidateToken implements domain.UserService.
func (u userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := u.cacheRepository.Get("user" + token)

	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	var user domain.User
	_ = json.Unmarshal(data, &user)

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

func (u userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	log.Printf("registering user: %s", req.Username)
	exist, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("error finding user: %s", err.Error())
		return dto.UserRegisterRes{}, err
	}

	if exist != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameTaken
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Printf("error hashing password: %s", err.Error())
		return dto.UserRegisterRes{}, err
	}
	req.Password = string(hashedPass)

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err = u.userRepository.Insert(ctx, &user)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	otpCode := util.GenerateRandomString(4)
	referenceId := util.GenerateRandomString(16)

	log.Printf("your OTP: %s", otpCode)

	_ = u.emailService.Send(user.Email, "OTP Code", "OTP anda adalah: "+otpCode)

	_ = u.cacheRepository.Set("otp"+referenceId, []byte(otpCode))
	_ = u.cacheRepository.Set("user-ref"+referenceId, []byte(user.Username))

	return dto.UserRegisterRes{
		ReferenceID: referenceId,
	}, nil

}

// ValidateOTP implements domain.UserService.
func (u userService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := u.cacheRepository.Get("otp" + req.ReferenceID)
	if err != nil {
		log.Printf("error finding user1: %s", err.Error())
		return domain.ErrAuthFailed
	}

	otp := string(val)
	if otp != req.OTP {
		log.Printf("invalid otp: %s", req.OTP)
		return domain.ErrOtpInvalid
	}

	val, err = u.cacheRepository.Get("user-ref" + req.ReferenceID)
	if err != nil {
		log.Printf("error finding user2: %s", err.Error())
		return domain.ErrOtpInvalid
	}

	user, err := u.userRepository.FindByUsername(ctx, string(val))
	if err != nil {
		log.Printf("error finding user3: %s", err.Error())
		return err
	}

	user.EmailVerifiedAt = time.Now()
	err = u.userRepository.Update(ctx, &user)
	if err != nil {
		log.Printf("error finding user4: %s", err.Error())
		return err
	}

	return nil
}
