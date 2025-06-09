package service

import "errors"

var (
	ErrPlayerNotFound = errors.New(MsgPlayerNotFound)
	ErrPlayerNotOwned = errors.New(MsgPlayerNotOwned)
	ErrUpdateNoFields = errors.New(MsgPlayerUpdateNoFields)
	ErrInvalidCountry = errors.New(MsgInvalidCountry)
)
