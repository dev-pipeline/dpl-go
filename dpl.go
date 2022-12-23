package main

import (
	"github.com/dev-pipeline/dpl-go/cmd"
	_ "github.com/dev-pipeline/dpl-go/internal/cmake"
)

func main() {
	cmd.Execute()
}
