package main

import (
  "os"
  "github.com/mikattack/twfconf"
)


/* 
 * Defines the options and flags that may be used when running the API.
 * 
 * The following flags are supported:
 * 
 *   "log-level"  - Logging level (default: WARN)
 *   "port"       - TCP port of service (default: 48200)
 *   "content"    - Content stub directory (default: /var/tmp/rest-stub)
 * 
 * These flags may also be passed as environmental variables:
 * 
 *   - LOG_LEVEL
 *   - PORT
 *   - CONTENT
 * 
 * Both types of parameters may be passed, but commandline flags will
 * always override environmental ones.
 */


const (
  DEFAULT_LOG_LEVEL string  = "INFO"
  DEFAULT_PORT string       = "48200"
  DEFAULT_CONTENT string    = "/var/tmp/rest-stub"
)

var config map[string]string


/* 
 * Convenience method to extract arguments from the environment and
 * the commandline.
 */
func getConfig() map[string]string {
  if (config == nil) {
    config = resolveConfig(os.Args[1:])
  }
  return config
}


/* 
 * Read arguments from the environment and from an argument collection
 * (typically "os.Args[1:]").  Arguments from the argument collection
 * override those from the environment.
 */
func resolveConfig(args []string) map[string]string {
  config := twfconf.NewArgConf(
    "rest-stub [args]",
    "Stub HTTP server to aid in developing REST services quickly")

  config.NewArg("log-level", "LOG_LEVEL", DEFAULT_LOG_LEVEL, "Logging level")
  config.NewArg("port", "PORT", DEFAULT_PORT, "TCP port of service")
  config.NewArg("content", "CONTENT", DEFAULT_CONTENT, "Directory for stub content files")

  return config.GetArgValues(args)
}
