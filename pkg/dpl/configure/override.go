package configure

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func loadOverride(overridePath string, modSet modifierSet) error {
	return loadModifierConfig(overridePath, modSet)
}

func applyOverrides(configRoot string, overrides []string, project dpl.Project) error {
	overridesDir := path.Join(configRoot, "overrides.d")
	for _, componentName := range project.ComponentNames() {
		modSet := newModifierSet()
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
		component, err := project.GetComponent(componentName)
		if err != nil {
			return err
		}
		err = applyComponentModifiers(component, modSet)
		if err != nil {
			return err
		}
	}
	return nil
}
