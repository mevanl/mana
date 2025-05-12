package manaerror

import "errors"

var (
	ErrChannelFetchFailed = errors.New("could not get channel")
)
