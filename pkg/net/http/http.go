// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
)

func Post(ctx context.Context, url string, contentType string, data []byte, requestInjectors ...func(req *http.Request) error) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("empty url")
	}
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Keep-Alive", "300")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for _, injector := range requestInjectors {
		if injector == nil {
			continue
		}
		if err := injector(req); err != nil {
			return nil, fmt.Errorf("request inject failed, %w", err)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http post %s, %w", url, err)
	}

	return resp, nil
}
