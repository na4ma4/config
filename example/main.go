package main

import (
	"fmt"
	"log"

	"github.com/na4ma4/config"
)

// test-project.toml contains:
// *****************************
// [server]
// address="127.0.0.1:8080"
// *****************************

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
