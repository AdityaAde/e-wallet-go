package api

import (
	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
	"adityaad.id/belajar-auth/internal/util"
	"github.com/gofiber/fiber/v2"
)

type transferApi struct {
	transactionService domain.TransactionService
	factorService      domain.FactorService
}

func NewTransfer(app *fiber.App,
	transactionService domain.TransactionService,
	authMid fiber.Handler,
	factorService domain.FactorService,
) {
	h := &transferApi{
		transactionService: transactionService,
		factorService:      factorService,
	}

	app.Post("transfer/inquiry", authMid, h.TransferInquiry)
	app.Post("transfer/execute", authMid, h.TransferExecute)
}

func (t transferApi) TransferInquiry(ctx *fiber.Ctx) error {
	var req dto.TransferInquiryReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: "Invalid body",
		})
	}

	inquiry, err := t.transactionService.TransferInquiry(ctx.Context(), req)

	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(inquiry)
}

func (t transferApi) TransferExecute(ctx *fiber.Ctx) error {

	var req dto.TransferExecuteReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: "Invalid body",
		})
	}

	user := ctx.Locals("x-user").(dto.UserData)

	if err := t.factorService.ValidatePIN(ctx.Context(), dto.ValidatePinReq{
		PIN:    req.PIN,
		UserID: user.ID,
	}); err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	err := t.transactionService.TransferExecute(ctx.Context(), req)

	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.SendStatus(200)
}
