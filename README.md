# na4ma4/config

![CI](https://github.com/na4ma4/config/workflows/CI/badge.svg)
[![GoDoc](https://godoc.org/github.com/na4ma4/config/src/jwt?status.svg)](https://godoc.org/github.com/na4ma4/config)

Go package for a thread-safe config interface

## Installation

```shell
go get -u github.com/na4ma4/config
```

## Example

```golang
package main

import (
    "fmt"
    "log"

    "github.com/na4ma4/config"
)

func main() {
    // Create config (example supplied viper)
    // Supplied ViperConf takes a project name, then a list of file names,
    // if no filenames are found, the last one is considered where you want the config to be saved.
    vcfg := config.NewViperConfig("test-project2", "artifacts/test-project.toml", "/tmp/test-project.toml", "test/test-project.toml")

    server := vcfg.GetString("server.address")

    fmt.Printf("Server: %s\n", server)

    err := vcfg.Save()
    if err != nil {
        log.Fatal(err)
    }
}
```
