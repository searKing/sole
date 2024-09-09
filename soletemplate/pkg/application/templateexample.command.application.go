// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package application

import (
	"context"

	"github.com/searKing/sole/soletemplate/pkg/domain/templateexample"
)

type TemplateExampleHandler struct {
	templateFactory templateexample.Factory
}

func NewTemplateExampleHandler(f templateexample.Factory) TemplateExampleHandler {
	return TemplateExampleHandler{templateFactory: f}
}

func (h TemplateExampleHandler) Handle(ctx context.Context, req *templateexample.ExampleRequest) (resp *templateexample.ExampleResponse, err error) {
	service, err := h.templateFactory.NewTemplateExample()
	if err != nil {
		return nil, err
	}
	return service.Example(ctx, req)
}
