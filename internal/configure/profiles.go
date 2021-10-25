package configure

import (
	"errors"
	"fmt"
	"path"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func loadProfiles(configRoot string, profiles []string) (configfile.ModifierSet, error) {
	profilesDir := path.Join(configRoot, "profiles.d")
	modSet := configfile.NewModifierSet()
	for _, profileName := range profiles {
		profilePath := path.Join(profilesDir, fmt.Sprintf("%v.conf", profileName))
		err := LoadModifierConfig(profilePath, modSet)
		if err != nil {
			return modSet, err
		}
	}
	return modSet, nil
}

func applyProfiles(modSet configfile.ModifierSet, project dpl.Project) error {
	for _, componentName := range project.Components() {
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
