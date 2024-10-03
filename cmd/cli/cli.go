package main

import (
	"fmt"

	m "github.com/sharithg/siphon/api"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {

	t := typescriptify.New()
	t.CreateInterface = true
	t.BackupDir = ""

	t.Add(m.Node{})
	t.Add(m.GithubRepo{})

	err := t.ConvertToFile("../client/src/@types/api.ts")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("OK")
}
