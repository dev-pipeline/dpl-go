package configure

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Flags struct {
	BuildDir         string
	BuildDirBasename string
	ConfigFile       string
	Overrides        []string
	Profiles         []string
	RootDir          string
}

const (
	cacheDirName  string = ".dpl"
	cacheFileName string = "build.cache"
)

func (f Flags) getBuildDir() string {
	if len(f.BuildDir) > 0 {
		return f.BuildDir
	}
	if len(f.Profiles) > 0 {
		return fmt.Sprintf("%v-%v", f.BuildDirBasename, strings.Join(f.Profiles, ","))
	}
	return f.BuildDirBasename
}

func getCachePath(f Flags) (string, string) {
	buildDir := f.getBuildDir()
	cacheDir := path.Join(buildDir, cacheDirName)
	cacheFile := path.Join(cacheDir, cacheFileName)

	return cacheDir, cacheFile
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
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = applyProfiles(modSet, project)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = applyOverrides(dplConfigDir, flags.Profiles, project)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defaultComponent, err := project.getDefaultComponent()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defaultComponent.SetValue("dpl.profiles", flags.Profiles)
	defaultComponent.SetValue("dpl.overrides", flags.Overrides)

	cacheDir, cacheFile := getCachePath(flags)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	outConfig, err := os.Create(cacheFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer outConfig.Close()
	err = project.Write(outConfig)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
