// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: date.proto

// 日期查询 API

package date

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *DateRequest) Validate() error {
	return nil
}
func (this *DateResponse) Validate() error {
	if this.RequestId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("RequestId", fmt.Errorf(`value '%v' must not be an empty string`, this.RequestId))
	}
	return nil
}