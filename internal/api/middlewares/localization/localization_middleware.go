package middleware

import (
	"net/http"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/localization"
)

type LocalizationMiddleware struct {
	localizationService *localization.Service
}

func New(c Configuration) *LocalizationMiddleware {
	return &LocalizationMiddleware{
		localizationService: c.LocalizationService,
	}
}

func (lm LocalizationMiddleware) Localize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langHeader := r.Header.Get("Accept-Language")
		localizer := lm.localizationService.GetLocalizer(langHeader)

		ctx := apicontext.WithLocalizer(r.Context(), localizer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
