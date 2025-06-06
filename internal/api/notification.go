package api

import (
	"context"
	"time"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
	"github.com/gofiber/fiber/v2"
)

type NotificationApi struct {
	notificationService domain.NotificationService
}

func NewNotification(app *fiber.App, authMid fiber.Handler, notificationService domain.NotificationService) {
	h := &NotificationApi{
		notificationService: notificationService,
	}

	app.Get("/notifications", authMid, h.GetUserNotification)

}

func (n NotificationApi) GetUserNotification(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	user := ctx.Locals("x-user").(dto.UserData)

	notifications, err := n.notificationService.FindByUser(c, user.ID)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(notifications)
}
