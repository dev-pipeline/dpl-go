package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dpl",
		Short: "Chain commands together to build projects",
		Run:   doDpl,
	}
)

func doDpl(cmd *cobra.Command, args []string) {
	fmt.Printf("Args %v\n", args)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
