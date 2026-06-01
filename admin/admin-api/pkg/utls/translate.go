package utls

import (
	"admin-api/pkg/translate"

	"github.com/gofiber/fiber/v3"
)

func Translate(MessageID string, param *string, c fiber.Ctx) string {
	var data map[string]any
	if param != nil {
		data = map[string]any{"name": param}
	}

	msg, err := translate.TranslateWithError(c, MessageID, data)
	if err != nil {
		return MessageID
	}
	return msg
}
