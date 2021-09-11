package main

import (
	"flag"
	"fmt"
	"os"
)

// parse input flags
func (cfg *config) parseConfig() {
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.JSONFile, "JSONFile", "db.json", "JSON File to serve")

	displayVersion := flag.Bool("version", false, "Display version and exit")
	displayFlags := flag.Bool("displayflags", false, "Display flags")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	if *displayFlags {
		cfg.printConfig()
	}

}

func (c *config) printConfig() {
	fmt.Println("port     :", c.port)
	fmt.Println("JSONFile :", c.JSONFile)
}
