package configure

import (
	"log"
	"os"
	"path"
)

type Flags struct {
	BuildDir         string
	BuildDirBasename string
	ConfigFile       string
	Overrides        []string
	Profiles         []string
	RootDir          string
}

func DoConfigure(flags Flags, args []string) {
	project, err := loadConfig(flags.ConfigFile)
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
	err = applyOverrides(dplConfigDir, flags.Profiles, project)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
