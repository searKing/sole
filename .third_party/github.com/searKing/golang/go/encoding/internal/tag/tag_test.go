// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/searKing/golang/go/encoding/internal/tag"
)

type inputType struct {
	Name        Name              `default:"Alice"`
	Age         int               `default:"10"`
	IntArray    []int             `default:"[1,2,3]"`
	StringArray []string          `default:"[\"stdout\",\"./logs\"]"`
	Map         map[string]string `default:"{\"name\": \"Alice\", \"age\": \"18\"}"`
}
type Name string

func (name *Name) TagDefault() error {
	if *name == "" {
		*name = "Bob"
	}
	return nil
}
func TestTag(t *testing.T) {
	i := &inputType{}
	expect := &inputType{
		Name:        "Bob",
		Age:         10,
		IntArray:    []int{1, 2, 3},
		StringArray: []string{"stdout", "./logs"},
		Map:         map[string]string{"name": "Alice", "age": "18"},
	}
	err := tag.Tag(i, func(val reflect.Value, tag reflect.StructTag) error {
		return json.Unmarshal([]byte(tag.Get("default")), val.Addr().Interface())
	})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, expect) {
		t.Errorf("expect\n[\n%v\n]\nactual[\n%v\n]", expect, i)
	}
}
