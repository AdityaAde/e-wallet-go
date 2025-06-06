package main

import (
	"adityaad.id/belajar-auth/internal/api"
	"adityaad.id/belajar-auth/internal/component"
	"adityaad.id/belajar-auth/internal/config"
	"adityaad.id/belajar-auth/internal/middleware"
	"adityaad.id/belajar-auth/internal/repository"
	"adityaad.id/belajar-auth/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := repository.NewRedisClient(cnf)

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationRepository)
	notificationService := service.NewNotificationService(notificationRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, transactionService, authMid)
	api.NewNotification(app, authMid, notificationService)

	_ = app.Listen((cnf.Server.Host + ":" + cnf.Server.Port))
}
