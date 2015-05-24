package main

import (
  "os"
  "reflect"
  "testing"
)



var envs []string = []string{ "DEBUG","LOG","PORT","CONTENT" }
var keys []string = []string{ "log","port","content" }


func resetEnv() {
  for _, env := range envs {
    os.Setenv(env, "")
  }
}


// Ensure default program options are what we expect
func TestConfigFlagDefaults(t *testing.T) {
  resetEnv()

  config, debug := resolveConfig(os.Args[1:])

  if (debug == true) {
    t.Error("[debug] expected to be [TRUE], found [FALSE]")
  }

  expected := map[string]string{
    "log":        "/var/log/rest-stub",
    "port":       "48200",
    "content":    "/var/tmp/rest-stub",
  }

  for key, value := range expected {
    if (config[key] != value) {
      t.Errorf("[%s] expected to be [%s], found [%s]", key, value, config[key])
    }
  }
}


// Test that environmental variables successfully override defaults
func TestConfigEnv(t *testing.T) {
  resetEnv()

  // Set environment
  for _, env := range envs {
    os.Setenv(env, "8")
  }
  os.Setenv("DEBUG", "true")

  // Test
  config, debug := resolveConfig(os.Args[1:])
  if (debug == false) {
    t.Error("[debug] expected to be [TRUE], found [FALSE]")
  }
  for _, key := range keys {
    if (config[key] != "8") {
      t.Errorf("[%s] expected to be [8], found [%s]", key, config[key])
    }
  }
}


// Test that commandline flags successfully override defaults
func TestConfigCLI(t *testing.T) {
  resetEnv()

  // Set commandline
  args := []string{}
  for _, key := range keys {
    args = append(args, "--" + key + "=2")
  }
  args = append(args, "--debug=true")

  // Test
  config, debug := resolveConfig(args)
  if (debug == false) {
    t.Error("[debug] expected to be [TRUE], found [FALSE]")
  }
  for _, key := range keys {
    if (config[key] != "2") {
      t.Errorf("[%s] expected to be [2], found [%s]", key, config[key])
    }
  }
}


// Test that commandline flags successfully override environmental ones
func TestConfigOverride(t *testing.T) {
  resetEnv()

  // Set environment
  for _, env := range envs {
    os.Setenv(env, "8")
  }
  os.Setenv("DEBUG", "true")

  // Set commandline
  args := []string{}
  for _, key := range keys {
    args = append(args, "--" + key + "=2")
  }
  args = append(args, "--debug=true")

  // Test
  config, debug := resolveConfig(args)
  if (debug == false) {
    t.Error("[debug] expected to be [TRUE], found [FALSE]")
  }
  for _, key := range keys {
    if (config[key] != "2") {
      t.Errorf("[%s] expected to be [2], found [%s]", key, config[key])
    }
  }
}

func TestConfigEquivalent(t *testing.T) {
  resetEnv()

  a, ad := getConfig()
  b, bd := resolveConfig(os.Args[1:])

  if reflect.DeepEqual(a, b) == false || ad != bd {
    t.Error("Convenience and parsing functions do not produce the same results")
  }
}
