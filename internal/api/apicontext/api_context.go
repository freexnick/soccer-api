package apicontext

import (
	"context"
	"soccer-api/internal/domain/entity"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AuthKey string
type languageKey string

const (
	AuthContextKey      AuthKey     = "authContext"
	LocalizerContextKey languageKey = "localizer"
)

func WithUser(ctx context.Context, userInfo *entity.Credentials) context.Context {
	return context.WithValue(ctx, AuthContextKey, userInfo)
}

func WithLocalizer(ctx context.Context, localizer *i18n.Localizer) context.Context {
	return context.WithValue(ctx, LocalizerContextKey, localizer)
}

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	val, _ := ctx.Value(LocalizerContextKey).(*i18n.Localizer)
	return val
}

func GetUser(ctx context.Context) *entity.Credentials {
	val, _ := ctx.Value(AuthContextKey).(*entity.Credentials)
	return val
}
