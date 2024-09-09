// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import "context"

// EventSourcedAggregator 领域事件聚合接口
// 类似领域事件的Handler
type EventSourcedAggregator interface {
	// apply 根据实际业务规则和更新状态来处理各种不同的领域事件，相当于领域事件的路由分发并调用处理
	// 该方法不应该被聚合外部所调用，聚合外部应该调用封装好的业务接口
	// 更新领域事件聚合 EventSourcedAggregate 的Version字段
	apply(ctx context.Context, changes DomainEvent)
}
