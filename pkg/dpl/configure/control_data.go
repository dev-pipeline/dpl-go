package configure

import (
	"path"
)

const (
	controlSectionName string = "dpl.control"

	subfilesKey string = "subfile"
)

type controlData struct {
	subfiles []string
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
	return nil
}
