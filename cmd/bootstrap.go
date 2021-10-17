package cmd

import (
	"github.com/spf13/cobra"
)

var (
	bootstrapCommon commonArgs

	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap a dpl project",
		Run:   doBootstrap,
	}
)

func doBootstrap(cmd *cobra.Command, args []string) {
}

func init() {
	addCommonArgs(bootstrapCmd, &bootstrapCommon)
	rootCmd.AddCommand(bootstrapCmd)
}
