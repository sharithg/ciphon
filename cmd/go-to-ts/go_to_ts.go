package main

import (
	"flag"
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

	if err != nil {
		log.Fatalf("error reading input files: %s", err)
	}

	var files []string

	for _, f := range pattens {
		files = append(files, filepath.Join(basepath, f))
	}

	parser := tools.NewGoToTs(files, prefix, tsPrefix, outputFile)

	if err = parser.ToTs(); err != nil {
		log.Fatalf("error converting to ts: %s\n", err)
	}
}

func main() {

	inputFiles := flag.String("input", "./**/*.go", "Input file pattern")
	outputFile := flag.String("output", "web/src/types/api.ts", "Output file path")
	prefix := flag.String("prefix", "Ts", "Prefix for generated types")
	tsPrefix := flag.String("tsPrefix", "T", "TypeScript prefix")
	flag.Parse()

	RunGoToTs(*inputFiles, *outputFile, *prefix, *tsPrefix)
}
