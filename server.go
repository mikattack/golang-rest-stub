package main

import (
  "bytes"
  "fmt"
  "log"
  "net/http"
  "strconv"
  "time"
)


func handleRequest(res http.ResponseWriter, req *http.Request) {
  var delay time.Duration = 0
  status := 200
  writer := bytes.Buffer{}

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

  // Use stub content for response body, when possible

  // Respond
  time.Sleep(delay)
  res.WriteHeader(status)
  if _, err := writer.WriteTo(res); err != nil {
    log.Printf("[%s]: Failed to write response\n", req.URL)
  }
}


func main() {
  config, _ := getConfig()
  fmt.Printf("Starting HTTP server on port [%s]\n", config["port"])
  http.HandleFunc("/", handleRequest)
  http.ListenAndServe(":" + config["port"], nil)
}
