// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeader(t *testing.T) {
	m := Header(func(header http.Header) {
		header.Set("foo", "bar")
	})
	h := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, nil)
	if "bar" != w.Header().Get("foo") {
		t.Errorf("expected %q, got %q", "bar", w.Header().Get("foo"))
	}
}
