package analyzer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"../base"
	dl "../downloader"
)

var analyzerIdgenerator dl.IdGenerator = dl.NewIdGenerator()

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]base.Data, []error)

type Analyzer interface {
	Id() uint32
	Analyze(respParsers []ParseResponse, resp base.Response) ([]base.Data, []error)
}

type AnalyzerPool interface {
	Take() (Analyzer, error)
	Return(analyzer Analyzer) error
	Total()uint32
	Used()uint32
}

type myAnalyzer struct {
	id uint32
}

func genAnalyzerId() uint32 {
	return analyzerIdgenerator.GetUint32()
}

func NewAnalyzer() Analyzer {
	return &myAnalyzer{id:genAnalyzerId()}
}

func (analyzer *myAnalyzer)Id() uint32 {
	return analyzer.id
}

func (analyzer *myAnalyzer)Analyze(respParsers []ParseResponse, resp base.Response) (dataList []base.Data, errorList []error) {
	if respParsers ==  nil {
		err := errors.New("The response parser list is invalid!")
		return nil, []error{err}
	}
	httpResp := resp.HttpResp()
	if httpResp == nil {
		err := errors.New("The http response is invalid")
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	log.Println(reqUrl)
	respDepth := resp.Depth()
	dataList = make([]base.Data, 0)
	errorList = make([]error, 0)
	for i, respParser := range respParsers {
		if respParser == nil {
			err := errors.New(fmt.Sprintf("The document parser[%d] is invalid!", i))
			errorList = append(errorList, err)
			continue
		}
		pDataList, pErrorList := respParser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				errorList = appendErrorList(errorList, pError)
			}
		}
	}
	return dataList, errorList
}

// 添加请求值或条目值到列表。
func appendDataList(dataList []base.Data, data base.Data, respDepth uint32) []base.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*base.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = base.NewRequest(req.HttpReq(), newDepth)
	}
	return append(dataList, req)
}

// 添加错误值到列表。
func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}
	return append(errorList, err)
}