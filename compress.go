// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// Compress is a compress middleware.
func Compress(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.CompressHandlerLevel(next, level)
	}
}
