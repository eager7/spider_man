package base

import (
	"bytes"
	"fmt"
)

type ErrorType string
type SpiderError interface {
	Type() ErrorType
	Error() string
}

type mySpiderError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

const (
	DOWNLOADER_ERROR     ErrorType = "Downloader Error"
	ANALYZER_ERROR       ErrorType = "Analyzer Error"
	ITEM_PROCESSOR_ERROR           = "Item Processor Error"
)

func NewSpiderError(errType ErrorType, errMsg string) SpiderError {
	return &mySpiderError{errType: errType, errMsg: errMsg}
}
func (s *mySpiderError) Type() ErrorType {
	return s.errType
}
func (s *mySpiderError) Error() string {
	if s.fullErrMsg == "" {
		s.getFullErrMsg()
	}
	return s.fullErrMsg
}
func (s *mySpiderError) getFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("Spider Man Error:")
	if s.errType != "" {
		buffer.WriteString(string(s.errType))
		buffer.WriteString(":")
	}
	buffer.WriteString(s.errMsg)
	s.fullErrMsg = fmt.Sprintf("%s\n", buffer.String())
	return
}
