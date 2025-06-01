package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type token struct {
	conf Configuration
}

func New(conf Configuration) repository.Token {
	return &token{conf: conf}
}

func (t *token) GenerateToken(ctx context.Context, user entity.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(t.conf.JWTExpiryMinutes) * time.Minute)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    t.conf.JWTIssuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.conf.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *token) extractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "bearer") {
		return "", errors.New("invalid authorization header format")
	}

	token := headerParts[1]
	if token == "" {
		return "", errors.New("token is missing")
	}

	return token, nil
}

func (t *token) ValidateToken(ctx context.Context, r *http.Request) (*entity.Token, error) {
	tokenString, err := t.extractTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	parsedClaims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(
		tokenString,
		parsedClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("token is malformed")
			}
			return []byte(t.conf.JWTSecret), nil
		},
		jwt.WithIssuer(t.conf.JWTIssuer),
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token is malformed")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token is expired or not yet valid")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("token signature is invalid")
		} else if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			return nil, errors.New("token issuer is invalid")
		}
		return nil, fmt.Errorf("jwt service: token parsing failed: %w", err)
	}

	if !jwtToken.Valid {
		return nil, errors.New("token is invalid")
	}

	return &entity.Token{
		UserID:    parsedClaims.UserID,
		Email:     parsedClaims.Email,
		Issuer:    parsedClaims.Issuer,
		IssuedAt:  parsedClaims.IssuedAt.Time,
		ExpiresAt: parsedClaims.ExpiresAt.Time,
	}, nil
}
