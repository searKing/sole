package date

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/go-grpc-middleware/interceptors/x_request_id"
	"github.com/searKing/sole/api/protobuf-spec/v1/date"
	"google.golang.org/grpc/metadata"
)

type Controller struct{}

// 日期查询
func (c *Controller) Now(ctx context.Context, req *date.DateRequest) (resp *date.DateResponse, err error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, errors.New("missing context value with grpc")
	}

	requestId := req.GetRequestId()
	if requestId == "" {
		requestIds := md.Get(x_request_id.DefaultXRequestIDKey)
		for _, id := range requestIds {
			if id != "" {
				requestId = id
				break
			}
		}
	}
	return &date.DateResponse{
		RequestId: requestId,
		Date:      time.Now().String(),
	}, nil
}
