package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	bootstrapCommon common.Args

	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap a dpl project",
		Run:   doBootstrap,
	}
)

func doBootstrap(cmd *cobra.Command, args []string) {
}

func init() {
	AddCommonArgs(bootstrapCmd, &bootstrapCommon)
	rootCmd.AddCommand(bootstrapCmd)
}
