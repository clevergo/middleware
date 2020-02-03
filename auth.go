// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
)

type contextKey int

const (
	basicAuthUserKey contextKey = iota
)

// GetBasicAuthUser returns the basic auth username.
func GetBasicAuthUser(r *http.Request) string {
	username, _ := r.Context().Value(basicAuthUserKey).(string)
	return username
}

// BasicAuthValidate is a function that validate username and password.
type BasicAuthValidate func(username, password string) bool

// BasicAuthOption is a function that apply on basic auth.
type BasicAuthOption func(*basicAuth)

// BasicAuthRealm is an option to sets realm.
func BasicAuthRealm(v string) BasicAuthOption {
	return func(ba *basicAuth) {
		ba.realm = v
	}
}

// BasicAuthErrorHandler is an option to handle validation failed.
func BasicAuthErrorHandler(h http.Handler) BasicAuthOption {
	return func(ba *basicAuth) {
		ba.errorHandler = h
	}
}

type basicAuth struct {
	realm        string
	validate     BasicAuthValidate
	errorHandler http.Handler
	next         http.Handler
}

func (ba *basicAuth) handleFailure(w http.ResponseWriter, r *http.Request) {
	if ba.errorHandler != nil {
		ba.errorHandler.ServeHTTP(w, r)
		return
	}

	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// ServeHTTP implements http.Handler.
func (ba *basicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+ba.realm+`"`)

	username, password, ok := r.BasicAuth()
	if !ok || !ba.validate(username, password) {
		ba.handleFailure(w, r)
		return
	}

	ctx := context.WithValue(r.Context(), basicAuthUserKey, username)
	r = r.WithContext(ctx)
	ba.next.ServeHTTP(w, r)
}

// BasicAuth returns a basic auth middleware.
func BasicAuth(validate BasicAuthValidate, opts ...BasicAuthOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return BasicAuthHandler(next, validate, opts...)
	}
}

// BasicAuthHandler returns a basic auth handler.
func BasicAuthHandler(next http.Handler, validate BasicAuthValidate, opts ...BasicAuthOption) http.Handler {
	ba := &basicAuth{
		validate: validate,
		realm:    "Restricted",
		next:     next,
	}

	for _, opt := range opts {
		opt(ba)
	}

	return ba
}
