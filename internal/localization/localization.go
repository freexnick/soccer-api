package localization

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"soccer-api/internal/domain/entity"
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

func (s *Service) LocalizePlayerPosition(localizer *i18n.Localizer, position entity.PlayerPosition) string {
	var messageID string

	switch position {
	case entity.Goalkeeper:
		messageID = LblPositionGoalkeeper
	case entity.Defender:
		messageID = LblPositionDefender
	case entity.Midfielder:
		messageID = LblPositionMidfielder
	case entity.Attacker:
		messageID = LblPositionAttacker
	default:
		messageID = string(position)
	}

	return s.GetMessage(localizer, messageID)
}

func (s *Service) LocalizeCountry(localizer *i18n.Localizer, country entity.Country) string {
	var messageID string

	switch country {
	case entity.GERMANY:
		messageID = LblCountryGermany
	case entity.BELGIUM:
		messageID = LblCountryBelgium
	case entity.FRANCE:
		messageID = LblCountryFrance
	case entity.PORTUGAL:
		messageID = LblCountryPortugal
	case entity.SPAIN:
		messageID = LblCountrySpain
	case entity.SCOTLAND:
		messageID = LblCountryScotland
	case entity.TURKEY:
		messageID = LblCountryTurkey
	case entity.AUSTRIA:
		messageID = LblCountryAustria
	case entity.ENGLAND:
		messageID = LblCountryEngland
	case entity.HUNGARY:
		messageID = LblCountryHungary
	case entity.SLOVAKIA:
		messageID = LblCountrySlovakia
	case entity.ALBANIA:
		messageID = LblCountryAlbania
	case entity.DENMARK:
		messageID = LblCountryDenmark
	case entity.NETHERLANDS:
		messageID = LblCountryNetherlands
	case entity.ROMANIA:
		messageID = LblCountryRomania
	case entity.SWITZERLAND:
		messageID = LblCountrySwitzerland
	case entity.SERBIA:
		messageID = LblCountrySerbia
	case entity.ITALY:
		messageID = LblCountryItaly
	case entity.CZECHIA:
		messageID = LblCountryCzechia
	case entity.SLOVENIA:
		messageID = LblCountrySlovenia
	case entity.CROATIA:
		messageID = LblCountryCroatia
	case entity.GEORGIA:
		messageID = LblCountryGeorgia
	case entity.UKRAINE:
		messageID = LblCountryUkraine
	case entity.POLAND:
		messageID = LblCountryPoland
	default:
		messageID = string(country)
	}

	return s.GetMessage(localizer, messageID)
}
