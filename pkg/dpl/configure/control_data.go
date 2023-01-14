package configure

import (
	"path"
)

const (
	controlSectionName string = "dpl.control"

	subfilesKey    string = "subfile"
	buildConfigKey string = "build_config"
	profilesKey    string = "profiles"
	overridesKey   string = "overrides"

	// these are put in each component, so prefix them with dpl
	sourceDirKey string = "dpl.source_dir"
	workDirKey   string = "dpl.work_dir"
)

type controlData struct {
	subfiles []string
	fields   map[string][]string
}

func getControlData(project *IniProject, sourceDirAbsPath string) (controlData, error) {
	controlSection, found := project.getAnyComponent(controlSectionName)
	if !found {
		return controlData{}, nil
	}
	subfiles, err := controlSection.ExpandValue(subfilesKey)
	if err != nil {
		return controlData{}, err
	}
	for i := range subfiles {
		subfiles[i] = path.Join(sourceDirAbsPath, subfiles[i])
	}
	return controlData{
		subfiles: subfiles,
		fields:   map[string][]string{},
	}, nil
}

func applyControlData(project *IniProject, cd controlData) error {
	section, err := project.config.NewSection(controlSectionName)
	if err != nil {
		return err
	}
	iniComponent := &IniComponent{
		config:  section,
		project: project,
	}
	iniComponent.SetValue(subfilesKey, cd.subfiles)
	for k, v := range cd.fields {
		iniComponent.SetValue(k, v)
	}
	return nil
}
