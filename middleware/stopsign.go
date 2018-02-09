package middleware

import (
	"fmt"
	"sync"
)

type StopSign interface {
	Sign() bool
	Signed() bool
	Reset()
	Deal(code string)
	DealCount() uint32
	DealTotal() uint32
	Summary() string
}

type myStopSign struct {
	signed bool
	rwmutex sync.RWMutex
	dealCountMap map[string]uint32
}

func NewStopSign() StopSign {
	ss := &myStopSign{
		dealCountMap:make(map[string]uint32),
	}
	return ss
}

func (ss *myStopSign) Sign() bool {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if ss.signed {
		return false
	}
	ss.signed = true
	return true
}

func (ss *myStopSign) Signed() bool {
	return ss.signed
}

func (ss *myStopSign) Deal (code string) {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if !ss.signed {
		return
	}
	if _, ok := ss.dealCountMap[code]; !ok {
		ss.dealCountMap[code] = 1
	} else {
		ss.dealCountMap[code] += 1
	}
}

func (ss *myStopSign) Reset() {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	ss.signed = false
	ss.dealCountMap = make(map[string]uint32)
}

func (ss *myStopSign)DealCount()uint32 {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	return uint32(len(ss.dealCountMap))
}

func (ss *myStopSign)DealTotal() uint32 {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	return uint32(len(ss.dealCountMap))
}

func (ss *myStopSign) Summary()string {
	return fmt.Sprintf("signed:%d, total:%d\n", ss.signed, len(ss.dealCountMap))
}