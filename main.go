package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/0verbyte/envmold/internal/mold"
)

const (
	appName    = "mold"
	appVersion = "v0.1.0"
)

var (
	moldTemplate string
	outputWriter string
	debug        bool
)

func init() {
	flag.StringVar(&moldTemplate, "template", "mold.yaml", "Path to the mold environment template file")
	flag.StringVar(&outputWriter, "output", "stdout", "Where environment variables will be written. File path or stdout")
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

func createMold(r io.Reader) (*mold.MoldTemplate, error) {
	m, err := mold.New(r)
	if err != nil {
		return nil, err
	}
	if err := m.Generate(); err != nil {
		return nil, err
	}
	return m, nil
}

func getMoldEnvironmentWriter(writerType string) mold.Writer {
	switch writerType {
	case mold.WriterStdout:
		fallthrough
	default:
		return &mold.StdoutWriter{}
	}
}

func main() {
	log.Printf("Running %s (%s)\n", appName, appVersion)

	f, err := os.Open(moldTemplate)
	if err != nil {
		log.Fatalf("Failed to open %s: %v\n", moldTemplate, err)
	}
	defer f.Close()

	moldData, err := createMold(f)
	if err != nil {
		log.Fatalln("Failed to create mold", err)
	}

	moldData.WriteEnvironment(getMoldEnvironmentWriter(outputWriter))
}
