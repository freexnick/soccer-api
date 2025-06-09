package service

import "errors"

var (
	ErrPriceRequired           = errors.New(MsgPriceRequired)
	ErrPlayerNotListed         = errors.New(MsgPlayerNotListed)
	ErrPlayerAlreadyListed     = errors.New(MsgPlayerAlreadyListed)
	ErrCannotBuyOwn            = errors.New(MsgCannotBuyOwn)
	ErrNotEnoughBudget         = errors.New(MsgNotEnoughBudget)
	ErrListingPermissionDenied = errors.New(MsgListingPermissionDenied)
)
