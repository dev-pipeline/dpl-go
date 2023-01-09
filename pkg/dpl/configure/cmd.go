package configure

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/cmd"
)

var (
	configureFlags   ConfigureFlags
	reconfigureFlags ReconfigureFlags

	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure a dpl project",
		Run: func(cmd *cobra.Command, args []string) {
			DoConfigure(configureFlags, args)
		},
	}

	reconfigureCmd = &cobra.Command{
		Use:   "reconfigure",
		Short: "Reconfigure an existing configuration",
		Run: func(cmd *cobra.Command, args []string) {
			DoReconfigure(reconfigureFlags, args)
		},
	}
)

func init() {
	configureCmd.PersistentFlags().StringVar(&configureFlags.BuildDir, "build-dir", "",
		"Directory to write build configuration. If specified, build-dir-prefix is ignored.")
	configureCmd.PersistentFlags().StringVar(&configureFlags.BuildDirBasename, "build-dir-prefix", "build",
		"Prefix to use for build directories")
	configureCmd.PersistentFlags().StringVar(&configureFlags.ConfigFile, "config-file", "build.config",
		"Project configuration file")
	configureCmd.PersistentFlags().StringSliceVar(&configureFlags.Overrides, "override", []string{},
		"Apply an override set")
	configureCmd.PersistentFlags().StringSliceVar(&configureFlags.Profiles, "profile", []string{},
		"Apply a profile")
	configureCmd.PersistentFlags().StringVar(&configureFlags.RootDir, "root-dir", "",
		"Root directory for source checkouts")
	cmd.AddCommand(configureCmd)

	reconfigureCmd.PersistentFlags().BoolVar(&reconfigureFlags.Append, "append", false,
		"Append new overrides/profiles instead of replacing")
	reconfigureCmd.PersistentFlags().StringSliceVar(&reconfigureFlags.Overrides, "override", []string{},
		"Apply an override set")
	reconfigureCmd.PersistentFlags().StringSliceVar(&reconfigureFlags.Profiles, "profile", []string{},
		"Apply a profile")
	cmd.AddCommand(reconfigureCmd)
}
