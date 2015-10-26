package main

import (
  "log"
  "os"
  "github.com/hashicorp/logutils"
)


/* 
 * All logging is done to STDOUT.
 * 
 * Logging levels are currently supported. A log entry's level is indicated
 * within the log message itself with a pattern like : "[LEVEL] message."
 * 
 * Examples:
 *  // Assuming minimum level is 'WARN'
 *  log.Print("[DEBUG] Debugging")  // Does not print
 *  log.Print("[WARN] Warning")     // Prints
 *  log.Print("[ERROR] Erring")     // Prints
 * 
 * Currently, the following levels are supported: DEBUG, INFO, WARN, ERROR.
 */
func init() {
  config := getConfig()
  filter := &logutils.LevelFilter{
    Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
    MinLevel: logutils.LogLevel(config["log-level"]),
    Writer:   os.Stdout,
  }
  log.SetOutput(filter)
  log.SetFlags(log.LstdFlags)
}