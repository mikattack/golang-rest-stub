package main // import "github.com/mikattack/golang-rest-stub"

import (
  "bufio"
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
  var (
    delay   time.Duration   = 0
    body    io.Reader       = req.Body
    output  io.Reader       = new(bytes.Buffer)
  )

  charset := "utf-8"
  mimetype := "application/octet-stream"
  status := 200

  withEcho := false
  content := "n/a"

  defer req.Body.Close()

  log.Printf("[DEBUG] Received request '%s'", req.URL)

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
    withEcho = true
  }

  // Use stub content for response body, if specified
  if val := req.Header.Get("X-Stub-Content"); val != "" {
    val = strings.Replace(val, "..", "", 0)
    fd, err := os.Open(config["content"] + "/" + val)
    if err != nil {
      fmt.Printf("[ERROR] Failed to read stub content: %s\n", config["content"] + "/" + val)
    } else {
      output = fd
      content = config["content"] + "/" + val
      withEcho = false
    }
    defer fd.Close()
  }

  // Buffer request body, so we can echo it back
  buffer := new(bytes.Buffer)
  if _, err := io.Copy(buffer, body); err != nil {
    log.Printf("[Error] Failed to buffer the request body\n")
  }
  reader := bytes.NewReader(buffer.Bytes())

  if withEcho == true {
    output = reader
  }

  // Echo request as debug output
  log.Printf("[DEBUG] Headers:\n")
  log.Printf("[DEBUG]   X-Stub-Delay: %s\n", delay)
  log.Printf("[DEBUG]   X-Stub-Status: %d\n", status)
  log.Printf("[DEBUG]   X-Stub-Content-Type: %s\n", mimetype)
  log.Printf("[DEBUG]   X-Stub-Charset: %s\n", charset)
  log.Printf("[DEBUG]   X-Stub-Echo: %t\n", withEcho)
  log.Printf("[DEBUG]   X-Stub-Content: %s\n", content)
  log.Printf("[DEBUG] Content (%d):\n", buffer.Len())

  scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
  for scanner.Scan() {
    log.Printf("[DEBUG]   %s\n", scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    log.Printf("[ERROR] Failed to read buffered request body\n")
  }

  // Respond
  time.Sleep(delay)
  res.Header().Set("Content-Type", mimetype + "; charset=" + charset)
  res.WriteHeader(status)
  if _, err := io.Copy(res, output); err != nil {
    log.Printf("[Error] Failed to write response\n")
    return
  }

  // Echo response as debug output
  log.Printf("[DEBUG] Sent %d response", status)
  log.Printf("[DEBUG] Headers:\n")
  for key, value := range res.Header() {
    log.Printf("[DEBUG]   %s: %s\n", key, value)
  }
  log.Printf("[DEBUG] Content:\n")
  scanner = bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
  for scanner.Scan() {
    log.Printf("[DEBUG]   %s\n", scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    log.Printf("[ERROR] Failed to read buffered response body\n")
  }
}


func main() {
  config := getConfig()
  fmt.Printf("[INFO] Starting HTTP server on port '%s'\n", config["port"])
  http.HandleFunc("/", handleRequest)
  http.ListenAndServe(":" + config["port"], nil)
}
