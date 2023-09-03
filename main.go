package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	appName    = "mold"
	appVersion = "v0.1.0"
)

var (
	environmentTemplate string
	environmentOutput   string
	debug               bool
)

func init() {
	flag.StringVar(&environmentTemplate, "template", "mold.yaml", "Path to the mold environment template file")
	flag.StringVar(&environmentOutput, "output", "stdout", "Where environment variables will be written. File path or stdout")
	flag.BoolVar(&debug, "debug", false, "Enables debug logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", appName, appVersion)
		flag.PrintDefaults()
	}
	flag.Parse()

	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func main() {
	log.Printf("Running %s (%s)\n", appName, appVersion)
}
