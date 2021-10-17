package cmd

import (
	"github.com/spf13/cobra"
)

type commonArgs struct {
	keepGoing    bool
	executor     string
	dependencies string
}

func addCommonArgs(command *cobra.Command, args *commonArgs) {
	command.PersistentFlags().BoolVar(&args.keepGoing, "keep-going", false,
		"Continue performing work even if a task fails")
	command.PersistentFlags().StringVar(&args.executor, "executor", "",
		"Method of executing work")
	command.PersistentFlags().StringVar(&args.dependencies, "dependencies", "",
		"Method of resolving dependencies")
}
