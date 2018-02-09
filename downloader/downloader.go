package downloader

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"sync"

	"../base"
	mdw "../middleware"
)

var downloaderGenerator IdGenerator = NewIdGenerator()

type PageDownloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response, error)
}

type IdGenerator interface {
	GetUint32() uint32
}

type PageDownloaderPool interface {
	Take() (PageDownloader, error)
	Return(dl PageDownloader) error
	Total() uint32
	Used() uint32
}

type myIdGenerator struct {
	sn uint32
	ended bool
	mutex sync.Mutex
}

type myPageDownloader struct {
	httpClient http.Client
	id uint32
}

func NewIdGenerator() IdGenerator {
	return &myIdGenerator{}
}

func (gen *myIdGenerator) GetUint32() uint32 {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	if gen.ended {
		defer func() {gen.ended = false}()
		gen.sn = 0
		return gen.sn
	}
	id := gen.sn
	if id < math.MaxUint32 {
		gen.sn ++
	} else {
		gen.ended = true
	}
	return id
}

func genDownloaderId() uint32 {
	return downloaderGenerator.GetUint32()
}

func NewPageDownloader(client *http.Client) PageDownloader {
	id := genDownloaderId()
	if client == nil {
		client = &http.Client{}
	}
	return &myPageDownloader{
		id:id,
		httpClient:*client,
	}
}

func (dl *myPageDownloader) Id() uint32 {
	return dl.id
}

func (dl *myPageDownloader) Download(req base.Request) (*base.Response, error) {
	httpReq := req.HttpReq()
	httpResp, err := dl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return base.NewResponse(httpResp, req.Depth()), nil
}

type myDownloaderPool struct {
	pool mdw.Pool
	etype reflect.Type
}

type GenPageDownloader func() PageDownloader

func NewPageDownloaderPool(total uint32, gen GenPageDownloader) (PageDownloaderPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() mdw.Entity {
		return gen()
	}
	pool ,err := mdw.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	dlpool := &myDownloaderPool{pool, etype}
	return dlpool, nil
}

func (dlpool *myDownloaderPool) Take() (PageDownloader, error) {
	entity, err := dlpool.pool.Take()
	if err != nil {
		return nil, err
	}
	dl, ok := entity.(PageDownloader)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is Not %s\n", dlpool.etype)
		panic(errors.New(errMsg))
	}
	return dl, nil
}

func (dlpool *myDownloaderPool) Return(dl PageDownloader) error {
	return dlpool.pool.Return(dl)
}

func (dlpool *myDownloaderPool) Total() uint32 {
	return dlpool.pool.Total()
}

func (dlpool *myDownloaderPool) Used() uint32 {
	return dlpool.pool.Used()
}