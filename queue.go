package main

import (
	"sync"
)

type UserType int64
const (
	UserTypeNone	UserType = iota
	//粉丝
	UserTypeFans	= 1 << 1
	//总督
	UserTypeZongDu	= 1 << 2
	//提督
	UserTypeTiDu	= 1 << 3
	//舰长
	UserTypeJianZhang	= 1 << 4
)

type QueueType	int
const(
	QueueTypeUnknow	QueueType = iota
	QueueTypeDanmu
	QueueTypeGift
	QueueTypeSpecialGift
)

type QueueFilter struct {
	mu					sync.Mutex
	UserEnable			UserType
	UserEnableLevel		int64
	QueueKind			QueueType
	QueueContent		string
}

var QueueFilterData	= defaultQueueFilter()

func addToQueue(data QueueItemData)  {
	listView.mu.Lock()
	defer listView.mu.Unlock()

	//去重
	listData := *listView.BindData
	for _, v := range listData {
		if data.Content == v.Content {
			return
		}
	}
	newList := append(*listView.BindData, data)
	listView.BindData = &newList
	if listView.List != nil {
		listView.List.Refresh()
	}
}

func defaultQueueFilter() QueueFilter {
	return QueueFilter{
		UserEnable:      0,
		UserEnableLevel: 0,
		QueueKind:       1,
		QueueContent:    "#排队",
	}
}