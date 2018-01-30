package analyzer

import (
	"net/http"

	"../base"
)

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]base.Data, []error)

type Analyzer interface {
	Id() uint32
	Analyze(respParsers []ParseResponse, resp base.Response) ([]base.Data, []error)
}

