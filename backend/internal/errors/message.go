package manaerror

import "errors"

var (
	ErrMessageEmpty       = errors.New("message is empty")
	ErrMessageSendFailed  = errors.New("failed to send message")
	ErrMessageFetchFailed = errors.New("failed to get message")
)
