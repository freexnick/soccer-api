package entity

import "github.com/google/uuid"

type Credentials struct {
	ID    uuid.UUID
	Email string
	Token string
}
