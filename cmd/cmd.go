package cmd

import (
	"github.com/spf13/cobra"

	icmd "github.com/dev-pipeline/dpl-go/internal/cmd"
)

func AddCommand(command *cobra.Command) error {
	return icmd.AddCommand(command)
}
