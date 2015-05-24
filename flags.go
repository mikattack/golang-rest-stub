package main

import (
  "os"
  "strconv"
  "github.com/spf13/cobra"
)


/* 
 * Defines the options and flags that may be used when running the API.
 * 
 * The following flags are supported:
 * 
 *   "debug"    - Output additional information to STDOUT.
 *   "log"      - Logfile directory (default: /var/log/rest-stub)
 *   "port"     - TCP port of service (default: 48200)
 *   "content"  - Content stub directory (default: /var/tmp/rest-stub)
 * 
 * These flags may also be passed as environmental variables:
 * 
 *   - DEBUG
 *   - LOG
 *   - PORT
 *   - CONTENT
 * 
 * Both types of parameters may be passed, but commandline flags will
 * always override environmental ones.
 */


const (
  DEFAULT_LOG string      = "/var/log/rest-stub"
  DEFAULT_PORT string     = "48200"
  DEFAULT_CONTENT string  = "/var/tmp/rest-stub"
)

var debug bool
var config map[string]string


/* 
 * Convenience method to extract arguments from the environment and
 * the commandline.
 */
func getConfig() (map[string]string, bool) {
  if (config == nil) {
    config, debug = resolveConfig(os.Args[1:])
  }
  return config, debug
}


/* 
 * Read arguments from the environment and from an argument collection
 * (typically "os.Args[1:]").  Arguments from the argument collection
 * override those from the environment.
 */
func resolveConfig(args []string) (map[string]string, bool) {
  var debug bool
  results := map[string]string{}
  flags := map[string]*string{}

  // Check environmental variables
  envs := map[string][]string{
    "log":      []string{"LOG", DEFAULT_LOG},
    "port":     []string{"PORT", DEFAULT_PORT},
    "content":  []string{"CONTENT", DEFAULT_CONTENT},
  }
  debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
  if err != nil {
    debug = false
  }
  for opt, definition := range envs {
    value := os.Getenv(definition[0])
    if len(value) > 0 {
      results[opt] = value
    } else {
      results[opt] = definition[1]
    }
  }

  // Override environmental variables with any specified commandline flags
  var RootCmd = &cobra.Command{
    Use:    "rest-stub",
    Short:  "REST API stubbing tool",
    Long:   "Simulate any REST API endpoint, for testing",
    Run:    func (cmd *cobra.Command, args []string) {
      // No operation, but necessary for Cobra to recognize a command
    },
  }

  RootCmd.SetArgs(args)   // Enables unit testing

  RootCmd.Flags().BoolVarP(&debug, "debug", "", debug, "Output additional information to STDOUT")
  flags["log"]      = RootCmd.Flags().String("log", results["log"], "Log file directory")
  flags["port"]     = RootCmd.Flags().String("port", results["port"], "TCP port of service")
  flags["content"]  = RootCmd.Flags().String("content", results["content"], "Directory for stub content files")
  
  RootCmd.Execute()

  for key, ptr := range flags {
    if len(*ptr) > 0 {
      results[key] = *ptr
    }
  }

  return results, debug
}
