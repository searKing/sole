// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package application

import (
	"github.com/google/wire"
)

// DefaultApplication is a wire provider set for the application that
// does not depend on the underlying platform.
var DefaultApplication = wire.NewSet(
	wire.Struct(new(Application), "*"),
	wire.Struct(new(Commands), "*"),
	wire.Struct(new(Queries), "*"),
	NewTemplateExampleHandler)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	TemplateExample TemplateExampleHandler
}

type Queries struct{}
