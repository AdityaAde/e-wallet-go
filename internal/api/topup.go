package api

import (
	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
	"github.com/gofiber/fiber/v2"
)

type topupApi struct {
	topupService domain.TopupService
}

func NewTopup(app *fiber.App, authMid fiber.Handler, topupService domain.TopupService) {
	t := topupApi{
		topupService: topupService,
	}

	app.Post("topup/initialize", authMid, t.InitializeTopup)
}

func (t topupApi) InitializeTopup(ctx *fiber.Ctx) error {
	var req dto.TopupReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-user").(dto.UserData)
	req.UserID = user.ID

	res, err := t.topupService.InitializeTopup(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(res)
}
