package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"adityaad.id/belajar-auth/dto"
	"github.com/gofiber/fiber/v2"
)

type notificationSse struct {
	hub *dto.Hub
}

func NewNotificationSse(app *fiber.App, authMid fiber.Handler, hub *dto.Hub) {
	h := &notificationSse{
		hub: hub,
	}

	app.Get("/sse/notification-stream", authMid, h.StreamNotification)
}

func (n notificationSse) StreamNotification(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/event-stream")

	user := ctx.Locals("x-user").(dto.UserData)
	n.hub.NotificationChannel[user.ID] = make(chan dto.NotificationData)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		event := fmt.Sprintf("event: %s\n"+"data: \n\n", "initial")
		_, _ = fmt.Fprintf(w, event)
		_ = w.Flush()

		log.Printf("sse connected: %d", user.ID)

		for notification := range n.hub.NotificationChannel[user.ID] {
			log.Printf("sse notification: %d", user.ID)
			data, _ := json.Marshal(notification)

			event := fmt.Sprintf("event: %s\n"+"data: %s\n\n", "notification-updated", data)

			_, _ = fmt.Fprint(w, event)
			_ = w.Flush()
		}
	})

	return nil
}
