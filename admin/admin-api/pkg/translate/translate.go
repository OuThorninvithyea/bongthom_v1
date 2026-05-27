package translate

import (

	// Commnuity pacakges
	"fmt"
	"log"
	"path/filepath"
	"github.com/gofiber/fiber/v3"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	// Internla pacakges
	"admin-api/pkg/logs"
	error_responses "admin-api/pkg/responses"
)

var bundle *goi18n.Bundle

func Init() *error_responses.ErrorResponse {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	localeFiles := []string{
		"pkg/i18n/localize/en.yaml",
		"pkg/i18n/localize/km.yaml",
		"pkg/i18n/localize/zh.yaml",
	}

	for _, file := range localeFiles {
		_, err := bundle.LoadMessageFile(filepath.Join(file))
		if err != nil {
			log.Printf("Error loading local file %s: %v", file, err)
			logs.NewCustomLog("translate_error", err.Error(), "error")
			return &error_responses.ErrorResponse{
				MessageID: "ErrorLoadMessage",
				Err:       err,
			}
		}
	}
	return nil
}

func TranslateWithError(c fiber.Ctx, key string, templateData ...map[string]any) (string, *error_responses.ErrorResponse) {
	if bundle == nil {
		logs.NewCustomLog("I18nNotInit", Init().ErrorString(), "error")
		return "", &error_responses.ErrorResponse{
			MessageID: key,
			Err:       fmt.Errorf("translation service is unavailable"),
		}
	}

	lang := c.Get("Accept-Language", "en")
	localizer := goi18n.NewLocalizer(bundle, lang)

	data := map[string]any{}
	if len(templateData) > 0 && templateData[0] != nil {
		data = templateData[0]
	}

	msg, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		log.Printf("Error localizing message ID %s: %v", key, err)
		logs.NewCustomLog("TranslationNotFound", err.Error(), "error")
		return "", &error_responses.ErrorResponse{
			MessageID: key,
			Err:       fmt.Errorf("Translation not found"),
		}
	}
	return msg, nil
}

func Translate(c fiber.Ctx, key string) string {
	msg, err := TranslateWithError(c, key)
	if err != nil {
		return key
	}
	return msg
}
