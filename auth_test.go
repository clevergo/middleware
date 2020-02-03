// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var testValidate = func(username, password string) bool {
	return username == "foo" && password == "bar"
}

func TestBasicAuth(t *testing.T) {
	handled := false
	username := ""
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
		username = GetBasicAuthUser(r)
	})
	h = BasicAuth(testValidate)(h)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	if handled {
		t.Error("basic auth failed")
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetBasicAuth("foo", "bar")
	h.ServeHTTP(w, r)
	if !handled {
		t.Error("basic auth failed")
	}
	if username != "foo" {
		t.Errorf("expected username %q, got %q", "foo", username)
	}
}

func TestBasicAuthRealm(t *testing.T) {
	tests := []string{"foo", "bar"}
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for _, realm := range tests {
		m := BasicAuthHandler(h, testValidate, BasicAuthRealm(realm))
		if ba, _ := m.(*basicAuth); realm != ba.realm {
			t.Errorf("expected realm %q, got %q", realm, ba.realm)
		}
	}
}

func TestBasicAuthErrorHandler(t *testing.T) {
	handled := false
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h = BasicAuthHandler(h, testValidate, BasicAuthErrorHandler(errorHandler))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !handled {
		t.Error("failed to set up error handler")
	}
}
