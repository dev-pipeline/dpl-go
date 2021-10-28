package configure

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func loadOverride(overridePath string, modSet configfile.ModifierSet) error {
	return LoadModifierConfig(overridePath, modSet)
}

func applyOverrides(configRoot string, overrides []string, project dpl.Project) error {
	overridesDir := path.Join(configRoot, "overrides.d")
	for _, componentName := range project.Components() {
		modSet := configfile.NewModifierSet()
		for _, override := range overrides {
			overridePath := path.Join(overridesDir, componentName, fmt.Sprintf("%v.conf", override))
			if _, err := os.Stat(overridePath); err == nil {
				err := loadOverride(overridePath, modSet)
				if err != nil {
					return err
				}
			} else if !errors.Is(err, os.ErrNotExist) {
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
