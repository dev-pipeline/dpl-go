package main

import (
	"github.com/dev-pipeline/dpl-go/cmd"
	_ "github.com/dev-pipeline/dpl-go/internal/cmake"
	_ "github.com/dev-pipeline/dpl-go/internal/git"
)

func main() {
	cmd.Execute()
}
