package configure

import (
	"fmt"
	"path"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func loadProfiles(configRoot string, profiles []string) (modifierSet, error) {
	profilesDir := path.Join(configRoot, "profiles.d")
	modSet := newModifierSet()
	for _, profileName := range profiles {
		profilePath := path.Join(profilesDir, fmt.Sprintf("%v.conf", profileName))
		err := loadModifierConfig(profilePath, modSet)
		if err != nil {
			return modSet, err
		}
	}
	return modSet, nil
}

func applyProfiles(modSet modifierSet, project dpl.Project) error {
	for _, componentName := range project.ComponentNames() {
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
