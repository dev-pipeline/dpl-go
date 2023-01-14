package cmd

import (
	"runtime"

	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/common"
)

func AddCommonArgs(command *cobra.Command, args *common.Args) {
	command.PersistentFlags().BoolVar(&args.KeepGoing, "keep-going", false,
		"Continue performing work even if a task fails")
	// command.PersistentFlags().StringVar(&args.Executor, "executor", "",
	// 	"Method of executing work")
	command.PersistentFlags().StringVar(&args.Dependencies, "dependencies", "deep",
		"Method of resolving dependencies")
	command.PersistentFlags().IntVar(&args.MaxTasks, "max-tasks", runtime.NumCPU(),
		"Maximum number of tasks to execute at once")
}
