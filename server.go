package main // import "github.com/mikattack/golang-rest-stub"

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "mime"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)


func handleRequest(res http.ResponseWriter, req *http.Request) {
  var delay time.Duration = 0
  var reader io.Reader = new(bytes.Buffer)
  charset := "utf-8"
  mimetype := "application/octet-stream"
  status := 200

  // Delay response, when requested
  if val := req.Header.Get("X-Stub-Delay"); val != "" {
    amount, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
      log.Printf("[%s]: Cannot interpret 'X-Stub-Delay' header value (%s) as integer\n", req.URL, val)
    } else {
      delay = time.Duration(amount) * time.Millisecond
    }
  }

  // Determine appropriate response status
  if val := req.Header.Get("X-Stub-Status"); val != "" {
    parsed, err := strconv.ParseInt(val, 10, 0)
    if err != nil || parsed < 100 || parsed >= 600 {
      log.Printf("[%s]: Invalid 'X-Stub-Status' value (%s)\n", req.URL, val)
    } else {
      status = int(parsed)
    }
  }

  // Determine appropriate MIME type
  if val := req.Header.Get("X-Stub-Content-Type"); val != "" {
    mediatype, _, err := mime.ParseMediaType(val)
    if err != nil {
      log.Printf("[%s]: Invalid MIME type specified (%s)\n", req.URL, val)
    } else {
      mimetype = mediatype
    }
  }

  // Determine character set (no validation here)
  if val := req.Header.Get("X-Stub-Charset"); val != "" {
    charset = val
  }

  // Echo request body, if specified
  // NOTE: The 'X-Stub-Content' header overrides this option
  if val := req.Header.Get("X-Stub-Echo"); val != "" {
    reader = req.Body
  }

  // Use stub content for response body, if specified
  if val := req.Header.Get("X-Stub-Content"); val != "" {
    val = strings.Replace(val, "..", "", 0)
    fd, err := os.Open(config["content"] + "/" + val)
    if err != nil {
      fmt.Printf("Failed to read stub content: %s\n", config["content"] + "/" + val)
    } else {
      reader = fd
    }
    defer fd.Close()
  }

  // Respond
  time.Sleep(delay)
  res.Header().Set("Content-Type", mimetype + "; charset=" + charset)
  res.WriteHeader(status)
  if _, err := io.Copy(res, reader); err != nil {
    log.Printf("[%s]: Failed to write response\n", req.URL)
  }
}


func main() {
  config, _ := getConfig()
  fmt.Printf("Starting HTTP server on port [%s]\n", config["port"])
  http.HandleFunc("/", handleRequest)
  http.ListenAndServe(":" + config["port"], nil)
}
