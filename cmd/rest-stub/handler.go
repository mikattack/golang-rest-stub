package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-kit/kit/log"
)

func RequestHandler(contentdir string, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			// Context keys defined by middleware
			requestId string = r.Context().Value(requestIdKey).(string)
			mime      string = r.Context().Value(mimeKey).(string)
			charset   string = r.Context().Value(charsetKey).(string)
			content   string = r.Context().Value(contentKey).(string)
		)

		// Determine appropriate response status
		status := 200
		if val := r.Header.Get("X-Stub-Status"); val != "" {
			parsed, err := strconv.ParseInt(val, 10, 0)
			if err != nil || parsed < 100 || parsed >= 600 {
				logger.Log(
					"warning", "invalid response status",
					"value", parsed,
					"requestId", requestId,
				)
				status = 400
			} else {
				status = int(parsed)
			}
		}

		// Load content based on middleware indicator
		var body *bytes.Buffer = new(bytes.Buffer)
		var err error = nil
		if content == "file" {
			file := r.Header.Get("X-Stub-Content")
			body, err = loadContent(file, contentdir)
		} else if content == "echo" {
			body, err = parseRequestBody(r)
		}
		if err != nil {
			logger.Log(
				"error", "failed to write content",
				"value", err,
				"requestId", requestId,
			)
			status = 500
		}

		// Send response
		w.WriteHeader(status)
		w.Header().Set("Content-Type", mime+"; charset="+charset)
		if _, err := io.Copy(w, body); err != nil {
			logger.Log(
				"error", "failed to write resposne",
				"value", err,
				"requestId", requestId,
			)
		}
	}
}

// Load contents of a file into a byte buffer from a file on disk in a
// given directory.
func loadContent(file string, dir string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	cleaned := filepath.Clean(filepath.Base(file))
	fullpath := fmt.Sprintf("%s/%s", dir, cleaned)
	fd, err := os.Open(fullpath)
	defer fd.Close()
	if err == nil {
		_, err = io.Copy(buffer, fd)
	}
	return buffer, err
}

// Load a Request's content into a byte buffer.
func parseRequestBody(r *http.Request) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, r.Body)
	return buffer, err
}
