package translation

import (
	"encoding/json"
	"errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Language struct {
	Tag      language.Tag
	FilePath string
}

type Translate struct {
	bundle *i18n.Bundle
}

func New(languages ...Language) (*Translate, error) {
	if len(languages) < 1 {
		return nil, errors.New("at least one language required")
	}

	bundle := i18n.NewBundle(languages[0].Tag)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lan := range languages {
		if _, lErr := bundle.LoadMessageFile(lan.FilePath); lErr != nil {
			return nil, lErr
		}
	}

	return &Translate{bundle: bundle}, nil
}

func (t Translate) TranslateMessage(messageID string, acceptLanguage ...string) string {
	localizer := i18n.NewLocalizer(t.bundle, acceptLanguage...)

	localizeConfig := i18n.LocalizeConfig{
		MessageID: messageID,
	}

	localization, lErr := localizer.Localize(&localizeConfig)
	if lErr != nil {
		return messageID
	}

	return localization
}
