package main

import (
	"flag"
	"log"

	"github.com/sharithg/siphon/api"
	"github.com/sharithg/siphon/internal/auth"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/tools"
)

func RunGoToTs(inputFiles, outputFile, prefix, tsPrefix string) {

	parser := tools.NewGoToTs(tsPrefix, outputFile)

	structs := []interface{}{
		api.Node{},
		repository.GetAllReposRow{},
		api.GithubRepoResponse{},
		api.ConnectRepoRequest{},
		repository.GetWorkflowRunsRow{},
		repository.GetWorkflowRunByIdRow{},
		repository.GetStepsByJobIdRow{},
		repository.GetCommandOutputsByStepIdRow{},
		repository.GetJobsByWorkflowIdRow{},
		auth.TokenPair{},
		repository.GetUserByIdRow{},
	}

	if err := parser.ToTs(structs); err != nil {
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
