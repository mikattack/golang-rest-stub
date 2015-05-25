package main

import (
  "fmt"
  "io"
  "log"
  "os"
)


/* 
 * Detects how the application should log based on runtime configuration
 * and sets the global logger appropriately.
 * 
 * By default, the logger logs to a file as set by the "log" configuration
 * directive.  If that directive is an empty string, then file logging is
 * disabled.
 * 
 * If the debugging has been enabled, the logger will additionally log
 * to STDOUT.  It is possible to enable debug output with no file logging.
 * 
 * Logging levels are not currently supported.
 */
func init() {
  config, debug := getConfig()
  writers := []io.Writer{}

  // Enable file logging (when possible)
  // NOT WORKING - Disabled for now
  if false && config["log"] != "" {
    fmt.Printf("Logging to file: %s\n", config["log"])
    fd, err := os.OpenFile(config["log"], os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
    if err != nil {
      fmt.Printf("File logging disabled, unable to access file\n")
      fmt.Println(err.Error())
    } else {
      writers = append(writers, fd)
    }
    defer fd.Close()
  }

  // Enable STDOUT debug logging (if requested)
  if debug {
    fmt.Println("Debugging enabled")
    writers = append(writers, os.Stdout)
  }

  log.SetOutput(io.MultiWriter(writers...))
  log.SetFlags(log.LstdFlags | log.Lshortfile)
}
