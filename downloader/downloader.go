package downloader
import "../base"

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
