// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import "net/http"

// HeaderFunc is a function that recieves a http.Header.
type HeaderFunc func(http.Header)

// Header is a HTTP midldeware that changes response header.
func Header(f HeaderFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f(w.Header())
			next.ServeHTTP(w, r)
		})
	}
}
