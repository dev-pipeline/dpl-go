package main

import (
	"github.com/dev-pipeline/dpl-go/internal/cmd"

	_ "github.com/dev-pipeline/dpl-go/plugins/cmake"
	_ "github.com/dev-pipeline/dpl-go/plugins/git"
)

func main() {
	cmd.Execute()
}
