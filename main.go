package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/0verbyte/envmold/internal/mold"
)

const (
	appName = "mold"
)

var (
	// set via ldflags at compile time.
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var (
	moldTemplate string
	outputWriter string
	debug        bool
	tags         string
)

func init() {
	flag.StringVar(&moldTemplate, "template", "mold.yaml", "Path to the mold environment template file")
	flag.StringVar(&outputWriter, "output", "stdout", "Where environment variables will be written. File path or stdout")
	flag.BoolVar(&debug, "debug", false, "Enables debug logging")
	flag.StringVar(&tags, "tags", "", "Filter environment variables matching tags")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", appName, appVersionString())
		flag.PrintDefaults()
	}
	flag.Parse()

	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func appVersionString() string {
	return fmt.Sprintf("%s-%s-%s", version, commit, date)
}

func createMold(r io.Reader) (*mold.MoldTemplate, error) {
	m, err := mold.New(r, getTags())
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

func getTags() *[]string {
	if tags == "" {
		return nil
	}
	s := strings.Split(tags, ",")
	return &s
}

func main() {
	log.Printf("Running %s (%s)\n", appName, appVersionString())

	f, err := os.Open(moldTemplate)
	if err != nil {
		log.Fatalf("Failed to open %s: %v\n", moldTemplate, err)
	}
	defer f.Close()

	moldData, err := createMold(f)
	if err != nil {
		log.Fatalln("Failed to create mold", err)
	}

	if err := moldData.WriteEnvironment(getMoldEnvironmentWriter(outputWriter)); err != nil {
		log.Fatalf("Failed to write mold to environment when using writer '%s': %v\n", outputWriter, err)
	}
}
