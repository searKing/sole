// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	errors_ "github.com/searKing/golang/go/errors"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/go/runtime"
	"github.com/searKing/sole/api/protobuf-spec/sole/types/v1/configuration"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
	"go.opentelemetry.io/otel"
)

func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", filepath.Base(os.Args[0])))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	runtime.LogPanic.AppendHandler(logPanic)
	runtime.NeverPanicButLog.AppendHandler(logPanic)

	errors_.ErrorHandlers = append(errors_.ErrorHandlers, func(err error) {
		if err != nil {
			caller, file, line := runtime.GetShortCallerFuncFileLine(3)
			slog.Error(fmt.Sprintf("Observed an error: %s at %s() %s:%d", err, caller, file, line))
		}
	})
}

type _log struct{}

func NewLog(ctx context.Context, cfg *v1.Configuration) (_ *_log, err error) {
	spanName := "NewLog"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := slog.With(logging.SpanAttrs(span)...)
	defer func() {
		if err != nil {
			logger.With(slog_.Error(err)).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	logs := cfg.GetLog()

	{
		var slogOpt slog.HandlerOptions
		slogOpt.Level = slog.Level(logs.GetLevel())
		slogOpt.AddSource = logs.GetAddSource()
		slogOpt.ReplaceAttr = slog_.ReplaceAttrTruncate(int(logs.GetTruncateAttrSizeTo()))
		var rotateOpts []slog_.RotateOption
		rotateOpts = append(rotateOpts,
			slog_.WithRotateRotateInterval(logs.GetRotationDuration().AsDuration()),
			slog_.WithRotateMaxCount(int(logs.GetRotationMaxCount())),
			slog_.WithRotateMaxAge(logs.GetRotationMaxAge().AsDuration()),
			slog_.WithRotateRotateSize(logs.GetRotationSizeInByte()))

		var newer func(path string, opts *slog.HandlerOptions, options ...slog_.RotateOption) (slog.Handler, error)
		switch logs.GetFormat() {
		case configuration.Log_json:
			newer = slog_.NewRotateJSONHandler
		case configuration.Log_text:
			newer = slog_.NewRotateTextHandler
		case configuration.Log_glog:
			newer = slog_.NewRotateGlogHandler
		case configuration.Log_glog_human:
			newer = slog_.NewRotateGlogHumanHandler
		default:
			// nop if unsupported
			return &_log{}, nil
		}
		var handlers []slog.Handler
		if logs.GetAllowStdout() {
			slogOpt2 := slogOpt
			slogOpt.Level = slog.Level(logs.GetLevel())
			h, err := newer("", &slogOpt2)
			if err != nil {
				slog.With("path", logs.GetPath()).
					With("duration", logs.GetRotationDuration().AsDuration()).
					With("max_count", logs.GetRotationMaxCount()).
					With("max_age", logs.GetRotationMaxAge().AsDuration()).
					With("rotate_size_in_byte", logs.GetRotationSizeInByte()).
					With("allow_stdout", logs.GetAllowStdout()).
					With("stdout_level", logs.GetStdoutLevel()).
					Info("add rotation wrapper for log")
				return nil, err
			}
			handlers = append(handlers, h)
		}
		{
			h, err := newer(logs.GetPath(), &slogOpt, rotateOpts...)
			if err != nil {
				return nil, err
			}
			handlers = append(handlers, h)
		}
		slog.SetDefault(slog.New(slog_.MultiHandler(handlers...)))

		slog.With("path", logs.GetPath()).
			With("duration", logs.GetRotationDuration().AsDuration()).
			With("max_count", logs.GetRotationMaxCount()).
			With("max_age", logs.GetRotationMaxAge().AsDuration()).
			With("rotate_size_in_byte", logs.GetRotationSizeInByte()).
			With("allow_stdout", logs.GetAllowStdout()).
			With("stdout_level", logs.GetStdoutLevel()).
			Info("add rotation wrapper for log")
	}
	return &_log{}, nil
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
		slog.Error(fmt.Sprintf("Observed a panic: %s\n%s", r, stacktrace))
	} else {
		slog.Error(fmt.Sprintf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace))
	}
}
