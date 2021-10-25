package configure

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

type Flags struct {
	BuildDir         string
	BuildDirBasename string
	ConfigFile       string
	Overrides        []string
	Profiles         []string
	RootDir          string
}

func loadOverride(componentName string, override string, modSet configfile.ModifierSet) error {
	return nil
}

func applyOverrides(overrides []string, project dpl.Project) error {
	for _, componentName := range project.Components() {
		modSet := configfile.NewModifierSet()
		for _, override := range overrides {
			err := loadOverride(componentName, override, modSet)
			if err != nil {
				return err
			}
		}
		component, found := project.GetComponent(componentName)
		if !found {
			return errors.New(fmt.Sprintf("Internal error; component %v not found", componentName))
		}
		err := configfile.ApplyComponentModifiers(component, modSet)
		if err != nil {
			return err
		}
	}
	return nil
}

func DoConfigure(flags Flags, args []string) {
	project, err := configfile.LoadProjectConfig(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	dplConfigDir := path.Join(homedir, ".dev-pipeline")

	modSet, err := loadProfiles(dplConfigDir, flags.Profiles)
	err = applyProfiles(modSet, project)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = applyOverrides(flags.Profiles, project)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
