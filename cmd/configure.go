package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/configure"
	"github.com/dev-pipeline/dpl-go/internal/reconfigure"
)

var (
	configureFlags   configure.Flags
	reconfigureFlags reconfigure.Flags

	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure a dpl project",
		Run: func(cmd *cobra.Command, args []string) {
			configure.DoConfigure(configureFlags, args)
		},
	}

	reconfigureCmd = &cobra.Command{
		Use:   "reconfigure",
		Short: "Reconfigure an existing configuration",
		Run: func(cmd *cobra.Command, args []string) {
			reconfigure.DoReconfigure(reconfigureFlags, args)
		},
	}
)

func init() {
	configureCmd.PersistentFlags().StringVar(&configureFlags.BuildDir, "build-dir", "",
		"Directory to write build configuration. If specified, build-dir-prefix is ignored.")
	configureCmd.PersistentFlags().StringVar(&configureFlags.BuildDirBasename, "build-dir-prefix", "build",
		"Prefix to use for build directories")
	configureCmd.PersistentFlags().StringVar(&configureFlags.ConfigFile, "config", "build.config",
		"Project configuration file")
	configureCmd.PersistentFlags().StringSliceVar(&configureFlags.Overrides, "override", []string{},
		"Apply an override set")
	configureCmd.PersistentFlags().StringSliceVar(&configureFlags.Profiles, "profile", []string{},
		"Apply a profile")
	configureCmd.PersistentFlags().StringVar(&configureFlags.RootDir, "root-dir", "",
		"Root directory for source checkouts")
	rootCmd.AddCommand(configureCmd)

	reconfigureCmd.PersistentFlags().BoolVar(&reconfigureFlags.Append, "append", false,
		"Append new overrides/profiles instead of replacing")
	reconfigureCmd.PersistentFlags().StringSliceVar(&reconfigureFlags.Overrides, "override", []string{},
		"Apply an override set")
	reconfigureCmd.PersistentFlags().StringSliceVar(&reconfigureFlags.Profiles, "profile", []string{},
		"Apply a profile")
	rootCmd.AddCommand(reconfigureCmd)
}
