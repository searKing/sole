// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import "fmt"

// EventWrapper 封装和代表一个独立的领域事件，其中包含与其所归属的流有关的所有时间的数据以及元数据
type EventWrapper struct {
	Id            string      // 事件ID,所属事件流中唯一， "{EventStreamId}-{EventNumber}"
	Event         DomainEvent // 事件元数据
	EventStreamId string      // 事件流ID
	EventNumber   int         // 事件版本号
}

func NewEventWrapper(event DomainEvent, eventNumber int, streamStateId string) EventWrapper {
	return EventWrapper{
		Id:            fmt.Sprintf("%s-%d", streamStateId, eventNumber),
		Event:         event,
		EventStreamId: streamStateId,
		EventNumber:   eventNumber,
	}
}
