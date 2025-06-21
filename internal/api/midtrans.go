package api

import (
	"adityaad.id/belajar-auth/domain"
	"github.com/gofiber/fiber/v2"
)

type midtransApi struct {
	midtransService domain.MidtransService
	topupService    domain.TopupService
}

func NewMidtrans(app *fiber.App, midtransService domain.MidtransService, topupService domain.TopupService) {
	m := midtransApi{
		midtransService: midtransService,
		topupService:    topupService,
	}

	app.Post("/midtrans/payment-callback", m.paymentHandlerNotification)
}

func (m midtransApi) paymentHandlerNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return ctx.SendStatus(400)
	}

	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return ctx.SendStatus(400)
	}

	success, _ := m.midtransService.VerifyPayment(ctx.Context(), notificationPayload)

	if success {
		_ = m.topupService.ConfirmedTopup(ctx.Context(), orderId)
		return ctx.SendStatus(200)
	}

	return ctx.SendStatus(400)
}
