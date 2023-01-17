package configure

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	errProjectOutOfDate error = fmt.Errorf("project is out of date")
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
	ret := controlData{
		fields: map[string][]string{},
	}
	controlSection, found := project.getAnyComponent(controlSectionName)
	if !found {
		return ret, nil
	}
	subfiles, err := controlSection.ExpandValue(subfilesKey)
	if err != nil {
		return ret, err
	}
	for i := range subfiles {
		subfiles[i] = path.Join(sourceDirAbsPath, subfiles[i])
	}
	ret.subfiles = subfiles
	return ret, nil
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

func setupFlags(cf *ConfigureFlags, controlComponent *IniComponent) error {
	type singleFieldMap struct {
		key   string
		field *string
	}
	singleMapper := []singleFieldMap{
		{key: buildConfigKey, field: &cf.ConfigFile},
		{key: workDirKey, field: &cf.BuildDir},
		{key: sourceDirKey, field: &cf.RootDir},
	}

	for i := range singleMapper {
		value, err := dpl.GetSingleComponentValue(controlComponent, singleMapper[i].key)
		if err != nil {
			return err
		}
		*singleMapper[i].field = value
	}

	type multipleFieldMap struct {
		key   string
		field *[]string
	}
	multipleMapper := []multipleFieldMap{
		{key: profilesKey, field: &cf.Profiles},
		{key: overridesKey, field: &cf.Overrides},
	}
	for i := range multipleMapper {
		values := controlComponent.GetValue(multipleMapper[i].key)
		*multipleMapper[i].field = values
	}
	return nil
}

func projectUpToDate(modTime time.Time, controlComponent *IniComponent) (ConfigureFlags, error) {
	cf := ConfigureFlags{}

	err := setupFlags(&cf, controlComponent)
	if err != nil {
		return cf, err
	}

	fileInfo, err := os.Stat(cf.ConfigFile)
	if err != nil {
		return cf, err
	}
	if modTime.Before(fileInfo.ModTime()) {
		return cf, errProjectOutOfDate
	}
	subfiles := controlComponent.GetValue(subfilesKey)
	for i := range subfiles {
		fileInfo, err := os.Stat(subfiles[i])
		if err != nil {
			return cf, err
		}
		if modTime.Before(fileInfo.ModTime()) {
			return cf, errProjectOutOfDate
		}
	}
	return cf, nil
}
