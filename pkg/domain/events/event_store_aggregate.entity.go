// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

// EventSourcedAggregate 领域事件聚合，一般建议组合EventSourcedAggregate，省去重新定义基础类型
// 类似领域事件的Handler
type EventSourcedAggregate struct {
	Changes        []DomainEvent // 未提交的事件的集合
	Version        int           // 该聚合版本序列号，用于乐观锁；与该聚合事件流的每次Apply调用而变化
	InitialVersion int           // 标志
}
