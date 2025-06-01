package auth

import "time"

type Configuration struct {
	JWTExpiryMinutes time.Duration
	JWTIssuer        string
	JWTSecret        string
}
