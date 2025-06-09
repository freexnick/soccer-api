package service

import "errors"

var (
	ErrTeamNotFound   = errors.New(MsgTeamNotFound)
	ErrUserHasNoTeam  = errors.New(MsgUserHasNoTeam)
	ErrUpdateNoFields = errors.New(MsgTeamUpdateNoFields)
	ErrInvalidCountry = errors.New(MsgInvalidCountry)
)
