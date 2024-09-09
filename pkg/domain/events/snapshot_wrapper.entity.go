// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import "time"

// Snapshot 代表包含每个聚合的元数据、文件、快照的单个文档
type Snapshot interface{}

// SnapshotWrapper 代表单个快照
type SnapshotWrapper struct {
	StreamName string    // 当前快照所归属的事件流的名称
	Snapshot   Snapshot  // 快照元数据
	Created    time.Time // 当前快照被保存的时间
}
