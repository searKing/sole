// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soletemplate

import (
	"errors"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/errorgoogle"
	"github.com/searKing/sole/pkg/web"
	"github.com/searKing/sole/soletemplate/pkg/domain/templateexample"
)

func ApiError(err error) error {
	if err == nil {
		return nil
	}
	err, _ = web.ErrorChain(TemplateImageError)(err, false)
	return err
}

func TemplateImageError(err error, handled bool) (error, bool) {
	if handled {
		return err, true
	}
	if errors.Is(err, templateexample.ErrMessageEmpty) {
		return errorgoogle.Errore(v1.SoleTemplateErrorEnum_INVALID_ARGUMENT_MESSAGE_EMPTY, err), true
	}
	return err, false
}
