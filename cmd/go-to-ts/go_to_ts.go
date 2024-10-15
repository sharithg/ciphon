package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/sharithg/siphon/internal/tools"
)

func RunGoToTs(inputFiles, outputFile, prefix, tsPrefix string) {

	basepath, pattern := doublestar.SplitPattern(inputFiles)
	fsys := os.DirFS(basepath)

	pattens, err := doublestar.Glob(fsys, pattern)

	var files []string

	for _, f := range pattens {
		files = append(files, filepath.Join(basepath, f))
	}

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("error reading input files: %s", err)
	}

	parser := tools.NewGoToTs(files, prefix, tsPrefix, outputFile)

	parser.ToTs()
}

func main() {
	RunGoToTs("./**/*.go", "web/src/types/api.ts", "Ts", "T")
}
