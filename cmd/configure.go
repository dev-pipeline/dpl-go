package cmd

import (
	"github.com/spf13/cobra"
)

var (
	append           bool
	buildDir         string
	buildDirBasename string
	configFile       string
	overrides        []string
	profiles         []string
	rootDir          string

	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure a dpl project",
		Run:   doConfigure,
	}

	reconfigureCmd = &cobra.Command{
		Use:   "reconfigure",
		Short: "Reconfigure an existing configuration",
		Run:   doConfigure,
	}
)

func doConfigure(cmd *cobra.Command, args []string) {
}

func init() {
	configureCmd.PersistentFlags().StringVar(&configFile, "build-dir", "",
		"Directory to write build configuration. If specified, build-dir-prefix is ignored.")
	configureCmd.PersistentFlags().StringVar(&configFile, "build-dir-prefix", "build",
		"Prefix to use for build directories")
	configureCmd.PersistentFlags().StringVar(&configFile, "config", "build.config",
		"Project configuration file")
	configureCmd.PersistentFlags().StringSliceVar(&overrides, "override", []string{},
		"Apply an override set")
	configureCmd.PersistentFlags().StringSliceVar(&profiles, "profile", []string{},
		"Apply a profile")
	configureCmd.PersistentFlags().StringVar(&configFile, "root-dir", "",
		"Root directory for source checkouts")
	rootCmd.AddCommand(configureCmd)

	reconfigureCmd.PersistentFlags().BoolVar(&append, "append", false,
		"Append new overrides/profiles instead of replacing")
	reconfigureCmd.PersistentFlags().StringSliceVar(&overrides, "override", []string{},
		"Apply an override set")
	reconfigureCmd.PersistentFlags().StringSliceVar(&profiles, "profile", []string{},
		"Apply a profile")
	rootCmd.AddCommand(reconfigureCmd)
}
