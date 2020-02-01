// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
)

// Logging is a logging middleware.
func Logging(out io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(out, next)
	}
}

// CombinedLogging is a combined logging middleware.
func CombinedLogging(out io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(out, next)
	}
}

// CustomerLogging is a customer logging middleware.
func CustomerLogging(out io.Writer, formatter handlers.LogFormatter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.CustomLoggingHandler(out, next, formatter)
	}
}
