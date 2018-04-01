package main

import (
	"context"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chilts/sid"
	"github.com/go-kit/kit/log"
)

const (
	_            = iota
	charsetKey   = iota
	contentKey   = iota
	mimeKey      = iota
	requestIdKey = iota
)

// ==========================================================================

// Standard middleware interface that takes an http.Handler, wraps it,
// and returns a new one.
type Middleware func(http.Handler) http.Handler

// Chain of middleware handlers (and shameless ripoff of the Alice package).
type middlewareChain struct {
	chain []Middleware
}

// Defines an immutable middleware chain.
func MiddlewareChain(middleware ...Middleware) middlewareChain {
	return middlewareChain{chain: middleware}
}

func (c middlewareChain) makeChain(handler http.Handler) http.Handler {
	for i := range c.chain {
		handler = c.chain[len(c.chain)-1-i](handler) // Reverse order
	}
	return handler
}

// Creates the middleware chain, using the given http.HandlerFunc as
// the root handler.
func (c middlewareChain) Then(handler http.HandlerFunc) http.Handler {
	return c.makeChain(handler)
}

// ==========================================================================

// Custom ResponseWriter type that records the information necessary to
// log a response.
type bufferedResponseWriter struct {
	http.ResponseWriter
	status  int
	written int
}

func (w *bufferedResponseWriter) Write(buffer []byte) (int, error) {
	n, err := w.ResponseWriter.Write(buffer)
	w.written = w.written + n
	return n, err
}

func (w *bufferedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// ==========================================================================

// Returns a request's internal identifier, from its Context. If such an
// identifier isn not found in the Context, the value "none" is returned.
func getRequestId(r *http.Request) string {
	id := "none"
	rid := r.Context().Value(requestIdKey)
	if rid != nil {
		id = rid.(string)
	}
	return id
}

// ==========================================================================

// Generates a request identifier to the request context. If an X-Request-Id
// header already exists, then its value will be used.
func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-ID")
			if rid == "" {
				rid = sid.Id()
			}
			ctx := context.WithValue(r.Context(), requestIdKey, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Delay response by the number of seconds specified in the X-Stub-Delay'
// header of the request. The default value is 0 (immediate).
func Delay(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if val := r.Header.Get("X-Stub-Delay"); val != "" {
				amount, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					logger.Log(
						"warning", "invalid delay value",
						"value", val,
						"requestId", getRequestId(r),
					)
				} else {
					time.Sleep(time.Duration(amount) * time.Millisecond)
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Add the response MIME type to the request context according to the
// X-Stub-Content-Type header. The default value is "application/octet-stream".
func MimeType(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			mimetype := "application/octet-stream"
			if val := r.Header.Get("X-Stub-Content-Type"); val != "" {
				mediatype, _, err := mime.ParseMediaType(val)
				if err != nil {
					logger.Log(
						"warning", "invalid MIME type/parameters",
						"value", mediatype,
						"requestId", getRequestId(r),
					)
					mediatype = "application/octect-stream"
				} else {
					mimetype = mediatype
				}
			}
			ctx := context.WithValue(r.Context(), mimeKey, mimetype)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Add the response MIME type to the request context according to the
// X-Stub-Charset header of the request. The default value is "utf-8".
// No validation is done on the requested character set identifier.
func Charset() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			charset := "utf-8"
			if val := r.Header.Get("X-Stub-Charset"); val != "" {
				charset = val
			}
			ctx := context.WithValue(r.Context(), charsetKey, charset)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Add an indicator to the request context for what kind of response body was
// requested. Only three values are possible: echo, file, or none (default).
// The indicator is set based on the presence/values of the X-Stub-Echo
// and X-Stub-Content headers.
//
// If the X-Stub-Echo header is present, the indicator is set to "echo".
// This means the content of the request will be used as the response content.
//
// If the X-Stub-Content header is present, the indicator is set to "file"
// and the response output is set to the contents of a file. The file is
// indicated by the header value in tandem with the program's content
// directory. If no matching file is found, no content is sent in the
// response.
//
// Note that if both X-Stub-Echo and X-Stub-Content headers are present in
// the request, the latter takes presidence.
func ResponseContent() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			content := "none"
			if val := r.Header.Get("X-Stub-Echo"); val != "" {
				content = "echo"
			}
			if val := r.Header.Get("X-Stub-Content"); val != "" {
				content = "file"
			}
			ctx := context.WithValue(r.Context(), contentKey, content)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Logs the details of the HTTP request. Note that the content of requests
// is not included.
func LogRequest(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			stack := []interface{}{
				"event", "request",
				"requestId", getRequestId(r),
				"url", r.URL,
				"protocol", r.Proto,
				"remote", r.RemoteAddr,
				"content", r.ContentLength,
			}
			stack = append(stack, flattenHeaders(r.Header)...)
			logger.Log(stack...)

			// Use a bufferedResponseWriter to record Response information
			// for later logging.
			writer := bufferedResponseWriter{w, 200, 0}
			next.ServeHTTP(&writer, r)

			// Defer a function to write the given Response details
			defer func(bw bufferedResponseWriter, requestId string) {
				stack := []interface{}{
					"requestId", requestId,
					"status", bw.status,
				}
				stack = append(stack, flattenHeaders(bw.Header())...)
				logger.Log(stack...)
			}(writer, getRequestId(r))
		}
		return http.HandlerFunc(fn)
	}
}

func flattenHeaders(headers http.Header) []interface{} {
	stack := []interface{}{}
	for name, headers := range headers {
		values := []string{}
		for _, h := range headers {
			values = append(values, h)
		}
		stack = append(stack, name)
		stack = append(stack, strings.Join(values, ","))
	}
	return stack
}
