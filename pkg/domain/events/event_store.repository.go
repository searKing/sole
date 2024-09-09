// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import (
	"context"
	"fmt"
)

// EventStoreRepository 事件存储
type EventStoreRepository interface {
	// CreateNewStream 创建事件流，事件流完成事件的版本控制和封装, DomainEvent - > EventWrapper
	CreateNewStream(ctx context.Context, streamName string, domainEvents ...DomainEvent) error
	// AppendEventsToStream 将领域事件附加到事件流上；将领域事件持久化为包含其所归属流的唯一的EventWrapper文档, DomainEvent - > EventWrapper
	// 需要验证最近存储的历史事件版本号是否匹配该待存储的当前事件的版本号，否则返回 OptimisticConcurrencyError
	AppendEventsToStream(ctx context.Context, streamName string, expectedVersion int, domainEvents ...DomainEvent) error
	// GetStream 查询事件流的事件列表，根据事件版本或ID进行查询，以便支持快照加载
	// fromVersion,toVersion 指定应该从事件流中提取回的最早和最晚事件
	GetStream(ctx context.Context, streamName string, fromVersion int, toVersion int) (domainEvents []DomainEvent, err error)
	// AddSnapshot 查询指定事件流的最新快照
	AddSnapshot(ctx context.Context, streamName string, snapshot interface{}) error
	// GetLatestSnapshot 查询指定事件流的最新快照
	GetLatestSnapshot(ctx context.Context, streamName string) (snapshot interface{}, err error)
}

// OptimisticConcurrencyError 并发更新事件流错误
type OptimisticConcurrencyError struct {
	GotVersion  int
	WantVersion int
}

func (e *OptimisticConcurrencyError) Error() string {
	return fmt.Sprintf("concurrent Event writes: got %d; want %d", e.GotVersion, e.WantVersion)
}
