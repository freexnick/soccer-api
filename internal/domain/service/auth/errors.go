package service

import "errors"

var (
	ErrUserAlreadyExists          = errors.New(MsgUserAlreadyExists)
	ErrInvalidCredentials         = errors.New(MsgInvalidCredentials)
	ErrSessionCreationFailed      = errors.New(MsgSessionCreationFailed)
	ErrUserDetailsRetrievalFailed = errors.New(MsgUserDetailsRetrievalFailed)
	ErrValidationFailed           = errors.New(MsgValidationFailed)
)
