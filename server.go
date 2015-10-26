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

  withEcho := "false"
  content := "n/a"

  log.Printf("[DEBUG] Request '%s'", req.URL)

  // Delay response, when requested
  if val := req.Header.Get("X-Stub-Delay"); val != "" {
    amount, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
      log.Printf("[WARN] Cannot interpret 'X-Stub-Delay' header value '%s' as integer (%s)\n", val, req.URL)
    } else {
      delay = time.Duration(amount) * time.Millisecond
    }
  }

  // Determine appropriate response status
  if val := req.Header.Get("X-Stub-Status"); val != "" {
    parsed, err := strconv.ParseInt(val, 10, 0)
    if err != nil || parsed < 100 || parsed >= 600 {
      log.Printf("[WARN] Invalid 'X-Stub-Status' value '%s' (%s)\n", val, req.URL)
    } else {
      status = int(parsed)
    }
  }

  // Determine appropriate MIME type
  if val := req.Header.Get("X-Stub-Content-Type"); val != "" {
    mediatype, _, err := mime.ParseMediaType(val)
    if err != nil {
      log.Printf("[WARN] Invalid MIME type specified '%s' (%s)\n", val, req.URL)
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
    withEcho = "true"
  }

  // Use stub content for response body, if specified
  if val := req.Header.Get("X-Stub-Content"); val != "" {
    val = strings.Replace(val, "..", "", 0)
    fd, err := os.Open(config["content"] + "/" + val)
    if err != nil {
      fmt.Printf("[ERROR] Failed to read stub content: %s\n", config["content"] + "/" + val)
    } else {
      reader = fd
      content = config["content"] + "/" + val
    }
    defer fd.Close()
  }

  log.Printf("[DEBUG] X-Stub-Delay: %s\n", delay)
  log.Printf("[DEBUG] X-Stub-Status: %d\n", status)
  log.Printf("[DEBUG] X-Stub-Content-Type: %s\n", mimetype)
  log.Printf("[DEBUG] X-Stub-Charset: %s\n", charset)
  log.Printf("[DEBUG] X-Stub-Echo: %s\n", withEcho)
  log.Printf("[DEBUG] X-Stub-Content: %s\n", content)


  // Respond
  time.Sleep(delay)
  res.Header().Set("Content-Type", mimetype + "; charset=" + charset)
  res.WriteHeader(status)
  if _, err := io.Copy(res, reader); err != nil {
    log.Printf("[Error] Failed to write response\n", req.URL)
  }
}


func main() {
  config := getConfig()
  fmt.Printf("[INFO] Starting HTTP server on port '%s'\n", config["port"])
  http.HandleFunc("/", handleRequest)
  http.ListenAndServe(":" + config["port"], nil)
}
