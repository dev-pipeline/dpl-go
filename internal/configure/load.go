package configure

import (
	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func applyConfig(config *ini.File) (*IniProject, error) {
	project := IniProject{
		config: config,
	}
	for _, component := range config.Sections() {
		if component.Name() != ini.DefaultSection {
			projectComponent := IniComponent{
				config: component,
			}
			err := dpl.ValidateComponent(&projectComponent)
			if err != nil {
				return nil, err
			}
		}
	}
	err := dpl.ValidateProject(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func loadRawConfig(data []byte) (*IniProject, error) {
	configFile, err := configfile.LoadRawConfig(data)
	if err != nil {
		return nil, err
	}
	project, err := applyConfig(configFile)
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
	project, err := applyConfig(configFile)
	if err != nil {
		return nil, err
	}
	return project, nil
}
