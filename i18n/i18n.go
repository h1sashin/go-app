package translator

import (
	"context"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type TranslatorKey string

const (
	TranslatorKeyLocalizer TranslatorKey = "localizer"
)

func NewTranslator() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, lang := range []string{"pl", "en"} {
		bundle.LoadMessageFile(fmt.Sprintf("i18n/locale/%s.yaml", lang))
	}

	return bundle
}

func InjectLocalizer(ctx context.Context, loc *i18n.Localizer) context.Context {
	return context.WithValue(ctx, TranslatorKeyLocalizer, loc)
}

func ExtractLocalizer(ctx context.Context) *i18n.Localizer {
	return ctx.Value(TranslatorKeyLocalizer).(*i18n.Localizer)
}
