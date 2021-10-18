package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/dpl"
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

type InvalidComponentNameError struct {
	Name string
}

func (icne *InvalidComponentNameError) Error() string {
	return fmt.Sprintf("Invalid name: %v", icne.Name)
}

func doConfigure(cmd *cobra.Command, args []string) {
}

func validateComponentName(component *dpl.Component) error {
	matched, err := regexp.Match("^([a-zA-Z](?:([-_])?[a-zA-Z0-9])+)+$", []byte(component.Name))

	if err != nil {
		return err
	}
	if !matched {
		return &InvalidComponentNameError{
			Name: component.Name,
		}
	}
	return nil
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

	dpl.RegisterComponentValidator("component-name", validateComponentName)
}
