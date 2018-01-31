package middleware

import (
	"errors"
	"fmt"
	"sync"

	"../base"
)

type ChannelManagerStatus uint8

const (
	CHANNEL_MANAGER_STATUS_UNINITALIZED ChannelManagerStatus = 0
	CHANNEL_MANAGER_STATUS_INITIALIZED  ChannelManagerStatus = 1
	CHANNEL_MANAGER_STATUS_CLOSED       ChannelManagerStatus = 2
)
const defaultChanLen uint = 1

var statusNameMap = map[ChannelManagerStatus]string {
	CHANNEL_MANAGER_STATUS_UNINITALIZED:"uninitialized",
	CHANNEL_MANAGER_STATUS_INITIALIZED:"initialized",
	CHANNEL_MANAGER_STATUS_CLOSED:"close",
}

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

type myChannelManager struct {
	rwmutex sync.RWMutex
	channelLen uint
	reqCh chan base.Request
	respCh chan base.Response
	itemCh chan base.Item
	errorCh chan error
	status ChannelManagerStatus
}

func NewChannelManager(channelLen uint) ChannelManager {
	if channelLen == 0{
		channelLen = defaultChanLen
	}
	cm := &myChannelManager{}
	cm.Init(channelLen, true)
	return cm
}

func (cm *myChannelManager)Init(channelLen uint, reset bool) bool {
	if channelLen == 0{
		panic(errors.New("The channel length is invalid!"))
	}
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if cm.status == CHANNEL_MANAGER_STATUS_INITIALIZED && !reset {
		return false
	}
	cm.channelLen = channelLen
	cm.reqCh = make(chan base.Request, channelLen)
	cm.respCh = make(chan base.Response, channelLen)
	cm.itemCh = make(chan base.Item, channelLen)
	cm.errorCh = make(chan error, channelLen)
	cm.status = CHANNEL_MANAGER_STATUS_INITIALIZED
	return true
}

func (cm *myChannelManager)Close() bool {
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if cm.status != CHANNEL_MANAGER_STATUS_INITIALIZED {
		return false
	}
	close(cm.reqCh)
	close(cm.respCh)
	close(cm.itemCh)
	close(cm.errorCh)
	cm.status = CHANNEL_MANAGER_STATUS_CLOSED
	return true
}

func (cm *myChannelManager)checkStatus() error {
	if cm.status == CHANNEL_MANAGER_STATUS_INITIALIZED {
		return nil
	}
	statusName, ok := statusNameMap[cm.status]
	if !ok {
		statusName = fmt.Sprintf("%d", cm.status)
	}
	errMsg := fmt.Sprintf("The undesirable status of channel manager:%s!\n", statusName)
	return errors.New(errMsg)
}

func (cm *myChannelManager)ReqChan() (chan base.Request, error) {
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if err := cm.checkStatus(); err != nil {
		return nil, err
	}
	return cm.reqCh, nil
}

func (cm *myChannelManager)RespChan()(chan base.Response, error) {
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if err := cm.checkStatus(); err != nil {
		return nil, err
	}
	return cm.respCh, nil
}

func (cm *myChannelManager) ItemChan()(chan base.Item, error)  {
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if err := cm.checkStatus(); err != nil {
		return nil, err
	}
	return cm.itemCh, nil
}

func (cm *myChannelManager) ErrorChan() (chan error, error) {
	cm.rwmutex.Lock()
	defer cm.rwmutex.Unlock()
	if err := cm.checkStatus(); err != nil {
		return nil, err
	}
	return cm.errorCh, nil
}

func (cm *myChannelManager)ChannelLen() uint {
	cm.rwmutex.RLock()
	defer cm.rwmutex.RUnlock()
	return cm.channelLen
}

func (cm *myChannelManager) Status() ChannelManagerStatus {
	cm.rwmutex.RLock()
	defer cm.rwmutex.RUnlock()
	return cm.status
}

var chanmanSummaryTemplate = "status:%s," + "requestChannel:%d/%d,"  + "responseChannel:%d/%d," + "itemChannel:%d/%d" + "errChannel:%d/%d"
func (cm *myChannelManager)Summary() string {
	cm.rwmutex.RLock()
	defer cm.rwmutex.RUnlock()
	summary := fmt.Sprintf(chanmanSummaryTemplate, statusNameMap[cm.status],
		len(cm.reqCh), cap(cm.reqCh), len(cm.respCh), cap(cm.respCh),
			len(cm.itemCh), cap(cm.itemCh), len(cm.errorCh), cap(cm.errorCh))
	return summary
}
