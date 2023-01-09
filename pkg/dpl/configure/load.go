package configure

import (
	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func validateProject(project *IniProject) error {
	for _, component := range project.config.Sections() {
		if component.Name() != ini.DefaultSection {
			projectComponent := IniComponent{
				config: component,
			}
			err := dpl.ValidateComponent(&projectComponent)
			if err != nil {
				return err
			}
		}
	}
	return dpl.ValidateProject(project)
}

func loadRawConfig(data []byte) (*IniProject, error) {
	configFile, err := configfile.LoadRawConfig(data)
	if err != nil {
		return nil, err
	}
	project := &IniProject{
		config: configFile,
	}
	err = validateProject(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func loadConfig(path string) (*IniProject, error) {
	configFile, err := configfile.LoadProjectConfig(path)
	if err != nil {
		return nil, err
	}
	project := &IniProject{
		config:     configFile,
		configFile: path,
	}
	return project, nil
}
