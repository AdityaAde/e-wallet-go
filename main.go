package main

import (
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/api"
	"adityaad.id/belajar-auth/internal/component"
	"adityaad.id/belajar-auth/internal/config"
	"adityaad.id/belajar-auth/internal/middleware"
	"adityaad.id/belajar-auth/internal/repository"
	"adityaad.id/belajar-auth/internal/service"
	"adityaad.id/belajar-auth/internal/sse"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := repository.NewRedisClient(cnf)
	hub := &dto.Hub{
		NotificationChannel: make(map[int64]chan dto.NotificationData),
	}

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationRepository, hub)
	notificationService := service.NewNotificationService(notificationRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, transactionService, authMid)
	api.NewNotification(app, authMid, notificationService)

	sse.NewNotificationSse(app, authMid, hub)

	_ = app.Listen((cnf.Server.Host + ":" + cnf.Server.Port))
}
