package main

import (
  "log"
  "os"
  "github.com/hashicorp/logutils"
)


/* 
 * All logging is done to STDOUT.
 * 
 * Logging levels are currently supported, but currently not configurable.
 * A log entry's level is indicated within the log message itself with
 * a pattern like this: [LEVEL].
 * 
 * Examples:
 *  // Assuming minimum level is 'WARN'
 *  log.Print("[DEBUG] Debugging")  // Does not print
 *  log.Print("[WARN] Warning")     // Prints
 *  log.Print("[ERROR] Erring")     // Prints
 * 
 * Currently, the following levels are supported: DEBUG, WARN, ERROR.
 */
func init() {
  filter := &logutils.LevelFilter{
    Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
    MinLevel: logutils.LogLevel("WARN"),
    Writer:   os.Stdout,
  }
  log.SetOutput(filter)
  log.SetFlags(log.LstdFlags)
}