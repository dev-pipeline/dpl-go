package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/checkout"
	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	args common.Args

	checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a dpl project",
		Run:   doCheckout,
	}
)

func doCheckout(cmd *cobra.Command, components []string) {
	common.DoCommand(components, args, []common.Task{
		checkout.Task,
	})
}

func init() {
	addCommonArgs(checkoutCmd, &args)
	rootCmd.AddCommand(checkoutCmd)
}
