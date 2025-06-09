package middleware

import (
	"net/http"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/domain/entity"
	tokenService "soccer-api/internal/domain/service/token"
	"soccer-api/internal/infrastructure/observer"
	localizationService "soccer-api/internal/localization"
)

type AuthMiddleware struct {
	tokenService        tokenService.Token
	localizationService localizationService.Service
	observer            *observer.Observer
}

func New(c Configuration) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService:        c.TokenService,
		observer:            c.Observer,
		localizationService: *c.LocalizationService,
	}
}

func (am AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		localizer := apicontext.GetLocalizer(ctx)

		token, err := am.tokenService.ValidateToken(ctx, r)
		if err != nil {
			helpers.UnauthorizedResponse(w, am.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError))
			return
		}

		identityForContext := entity.Credentials{
			ID:    token.UserID,
			Email: token.Email,
		}

		ctx = apicontext.WithUser(r.Context(), &identityForContext)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
