// Package middleware provides HTTP middleware for the Creator API.
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Recoverer is a middleware that recovers from panics, logs the error with stack trace,
// and returns a 500 Internal Server Error response. This prevents the server from crashing
// due to unhandled panics in handlers.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				// Get request ID for correlation
				requestID := middleware.GetReqID(r.Context())

				// Capture stack trace for debugging
				stack := debug.Stack()

				// Log panic with full context and stack trace
				log.Error().
					Str("request_id", requestID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("remote_addr", r.RemoteAddr).
					Interface("panic", rvr).
					Str("stack_trace", string(stack)).
					Msg("Panic recovered")

				// Return 500 Internal Server Error to client
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// PrintPanicStack is a helper function to format panic information.
// Used internally by Recoverer middleware.
func PrintPanicStack(rvr interface{}) string {
	return fmt.Sprintf("PANIC: %v\n%s", rvr, debug.Stack())
}
