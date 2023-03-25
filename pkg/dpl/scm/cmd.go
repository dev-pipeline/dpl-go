package scm

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/cmd"
	icmd "github.com/dev-pipeline/dpl-go/internal/cmd"
	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	args common.Args

	checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a dpl project",
		RunE:  doCheckout,
	}
)

func doCheckout(cmd *cobra.Command, components []string) error {
	return common.DoCommand(components, args, []common.Task{
		CheckoutTask,
	})
}

func init() {
	icmd.AddCommonArgs(checkoutCmd, &args)
	cmd.AddCommand(checkoutCmd)
}
