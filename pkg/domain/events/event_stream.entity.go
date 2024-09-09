// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

// EventStream 代表一个事件流，仅包含元数据
type EventStream struct {
	Id      string // 事件流ID，aggregate type + id
	Version int    // 领域事件版本生成器，用于保证事件在本流中的唯一性
}

func NewEventStream(id string) *EventStream {
	return &EventStream{
		Id:      id,
		Version: 0,
	}
}

// RegisterEvent 注册领域事件至本事件流，封装并生成绑定本事件流的领域事件
func (s *EventStream) RegisterEvent(event DomainEvent) EventWrapper {
	s.Version++
	return NewEventWrapper(event, s.Version, s.Id)
}
