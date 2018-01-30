package middleware

import "../base"

type ChannelManagerStatus uint8

const (
	CHANNEL_MANAGER_STATUS_UNINITALIZED ChannelManagerStatus = 0
	CHANNEL_MANAGER_STATUS_INITIALIZED  ChannelManagerStatus = 1
	CHANNEL_MANAGER_STATUS_CLOSED       ChannelManagerStatus = 2
)

type ChannelManager interface {
	Init(channelLen uint, reset bool) bool
	Close() bool
	ReqChan() (chan base.Request, error)
	RespChan() (chan base.Response, error)
	ItemChan() (chan base.Item, error)
	ErrorChan() (chan error, error)
	ChannelLen() uint
	Status() ChannelManagerStatus
	Summary() string
}
