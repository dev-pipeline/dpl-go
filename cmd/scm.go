package cmd

import (
	"github.com/spf13/cobra"

	_ "github.com/dev-pipeline/dpl-go/internal/checkout"
)

var (
	checkoutCommon commonArgs

	checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a dpl project",
		Run:   doCheckout,
	}
)

func doCheckout(cmd *cobra.Command, args []string) {
}

func init() {
	addCommonArgs(checkoutCmd, &checkoutCommon)
	rootCmd.AddCommand(checkoutCmd)
}
