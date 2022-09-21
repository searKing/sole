// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/runtime"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", os.Args[0]))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	runtime.LogPanic.AppendHandler(logPanic)
	runtime.NeverPanicButLog.AppendHandler(logPanic)

	errors.ErrorHandlers = append(errors.ErrorHandlers, func(err error) {
		if err != nil {
			caller, file, line := runtime.GetShortCallerFuncFileLine(2)
			logrus.Errorf("Observed an error: %s at %s() %s:%d", err, caller, file, line)
		}
	})
}

type _log struct{}

func NewLog(ctx context.Context, cfg *configpb.Configuration) (_ *_log, err error) {
	spanName := "NewLog"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).
		WithField("span_id", span.SpanContext().SpanID())
	defer func() {
		if err != nil {
			logger.WithError(err).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	logs := cfg.GetLog()

	var fc logrus_.FactoryConfig
	fc.SetDefaults()
	fc.Level = logrus.Level(logs.GetLevel())
	fc.Format = logrus_.Format(logs.GetFormat())
	fc.Path = logs.GetPath()

	fc.RotationDuration = logs.GetRotationDuration().AsDuration()
	fc.RotationSizeInByte = logs.GetRotationSizeInByte()
	fc.RotationMaxCount = int(logs.GetRotationMaxCount())
	fc.RotationMaxAge = logs.GetRotationMaxAge().AsDuration()
	fc.ReportCaller = logs.GetReportCaller()
	fc.MuteDirectlyOutput = logs.GetMuteDirectlyOutput()
	fc.MuteDirectlyOutputLevel = logrus.Level(logs.GetMuteDirectlyOutputLevel())
	fc.TruncateMessageSizeTo = int(logs.GetTruncateMessageSizeTo())
	fc.TruncateKeySizeTo = int(logs.GetTruncateKeySizeTo())
	fc.TruncateValueSizeTo = int(logs.GetTruncateValueSizeTo())
	return &_log{}, logrus_.NewFactory(fc).Apply()
}

// logPanic logs the caller tree when a panic occurs (except in the special case of http.ErrAbortHandler).
func logPanic(r interface{}) {
	if r == nil || r == http.ErrAbortHandler {
		// honor the http.ErrAbortHandler sentinel panic value:
		//   ErrAbortHandler is a sentinel panic value to abort a handler.
		//   While any panic from ServeHTTP aborts the response to the client,
		//   panicking with ErrAbortHandler also suppresses logging of a stack trace to the server's error log.
		return
	}

	const size = 64 << 10
	stacktrace := runtime.GetCallStack(size)
	if _, ok := r.(string); ok {
		logrus.Errorf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		logrus.Errorf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
