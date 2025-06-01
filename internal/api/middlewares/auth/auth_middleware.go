package middleware

import (
	"net/http"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/domain/entity"
	tokenService "soccer-api/internal/domain/service/token"
	"soccer-api/internal/infrastructure/observer"
)

type AuthMiddleware struct {
	tokenService tokenService.Token
	observer     *observer.Observer
}

func New(c Configuration) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: c.TokenService,
		observer:     c.Observer,
	}
}

func (am AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := am.tokenService.ValidateToken(ctx, r)
		if err != nil {
			helpers.UnauthorizedResponse(w, err.Error())
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
