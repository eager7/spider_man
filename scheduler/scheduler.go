package scheduler

import (
	"net/http"

	"../analyzer"
)
import "../itempipeline"

type GenHttpClient func() *http.Client

type SchedSummary interface {
	String() string
	Detail() string
	Same(other SchedSummary)bool
}

type Scheduler interface {
	Start(channelLen uint,
		poolSize uint32,
		depth uint32,
		httpClientGenerator GenHttpClient,
		respParsers []analyzer.ParseResponse,
		itemProcessors itempipeline.ProcessItem,
		firstHttpReq *http.Request) (err error)
	Stop() bool
	Running() bool
	ErrorChan() <-chan error
	Idle() bool
	Summary(prefix string) SchedSummary
}
