package localization

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

type Service struct {
	bundle         *i18n.Bundle
	defaultLangTag language.Tag
}

func New() (*Service, error) {
	defaultLang := language.English
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	err := fs.WalkDir(localesFS, "locales", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		if _, err := bundle.LoadMessageFileFS(localesFS, path); err != nil {
			return fmt.Errorf("failed to load locale %q: %w", path, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("localization init failed: %w", err)
	}

	return &Service{bundle: bundle, defaultLangTag: defaultLang}, nil
}

func (s *Service) GetLocalizer(langHeader string) *i18n.Localizer {
	defaultLang := s.defaultLangTag.String()

	if langHeader == "" {
		return i18n.NewLocalizer(s.bundle, defaultLang)
	}

	preferredTags, _, err := language.ParseAcceptLanguage(langHeader)
	if err != nil || len(preferredTags) == 0 {
		return i18n.NewLocalizer(s.bundle, defaultLang)
	}

	for _, clientTag := range preferredTags {
		clientBase, _ := clientTag.Base()
		for _, supported := range s.bundle.LanguageTags() {
			supportedBase, _ := supported.Base()
			if clientBase == supportedBase {
				return i18n.NewLocalizer(s.bundle, clientTag.String())
			}
		}
	}

	return i18n.NewLocalizer(s.bundle, defaultLang)
}

func (s *Service) GetMessage(localizer *i18n.Localizer, messageID string, templateData ...any) string {
	config := &i18n.LocalizeConfig{MessageID: messageID}
	if len(templateData) > 0 && templateData[0] != nil {
		config.TemplateData = templateData[0]
	}

	msg, err := localizer.Localize(config)
	if err != nil {
		return fmt.Sprintf("[%s - TranslationMissingOrError]", messageID)
	}
	return msg
}
