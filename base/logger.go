package base

import (
	"log"

	"github.com/eager7/go/mlog"
)

// 创建日志记录器。
func NewLogger() *log.Logger {
	return mlog.Info
}
