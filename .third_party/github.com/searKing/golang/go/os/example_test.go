// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"time"
)
import os_ "github.com/searKing/golang/go/os"

func ExampleNewRotateFile() {
	file := os_.NewRotateFile("log/test.2006-01-02-15-04-05.log")
	defer file.Close()
	file.MaxCount = 5
	file.RotateInterval = 5 * time.Second
	file.MaxAge = time.Hour
	file.FileLinkPath = "log/s.log"
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		file.WriteString(time.Now().String())
		//if err := file.Rotate(false); err != nil {
		//	fmt.Printf("%d, err: %v", i, err)
		//}
	}
	// Output:
}

func ExampleNewRotateFileWithStrftime() {
	file := os_.NewRotateFileWithStrftime("log/test.%Y-%m-%d-%H-%M-%S.log")
	file.MaxCount = 5
	file.RotateInterval = 5 * time.Second
	file.MaxAge = time.Hour
	file.FileLinkPath = "log/s.log"
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		file.WriteString(time.Now().String())
		//if err := file.Rotate(false); err != nil {
		//	fmt.Printf("%d, err: %v", i, err)
		//}
	}
	// Output:
}

func ExampleDiskUsage() {
	total, free, avail, inodes, inodesFree, err := os_.DiskUsage("/tmp")
	if err != nil {
		return
	}

	_, _, _, _, _ = total, free, avail, inodes, inodesFree
	//fmt.Printf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d", total, free, avail, inodes, inodesFree)
	// total :499963174912 B, free: 57534603264 B, avail: 57534603264 B, inodes: 566386444, inodesFree: 561861360

	// Output:
}
