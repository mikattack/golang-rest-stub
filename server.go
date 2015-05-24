package main

import (
  "fmt"
  "log"
  "net/http"
)


/* 
 * Start the HTTP server.
 */
func main() {
  config, _ := getConfig()
  router := httprouter.New()

  fmt.Printf("Starting HTTP server on port [%s]\n", config["port"])

  // Start a server to listen for JSON-encoded checkmk queries via HTTP
  fmt.Printf("Starting HTTP server on port [%s]\n", config["port"])
  http.HandleFunc("/", newStubHandler(config))
  http.ListenAndServe(":" + config["port"], nil)
}
