package configure

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ConfigureFlags struct {
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

	loaderString string = "configure"
)

var (
	errMissingSourceDir     error = fmt.Errorf("project missing source directory information")
	errMissingWorkDir       error = fmt.Errorf("project missing work directory information")
	errCouldntLoadComponent error = fmt.Errorf("couldn't load component")
	errNoControlData        error = fmt.Errorf("no control data")
)

func (f ConfigureFlags) getBuildDir() string {
	if len(f.BuildDir) > 0 {
		return f.BuildDir
	}
	if len(f.Profiles) > 0 {
		return fmt.Sprintf("%v-%v", f.BuildDirBasename, strings.Join(f.Profiles, ","))
	}
	return f.BuildDirBasename
}

func getCachePath(f ConfigureFlags) (string, string) {
	buildDir := f.getBuildDir()
	cacheDir := path.Join(buildDir, cacheDirName)
	cacheFile := path.Join(cacheDir, cacheFileName)

	return cacheDir, cacheFile
}

func DoConfigure(flags ConfigureFlags, args []string) {
	project, err := loadConfig(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	sourceFileAbsPath, err := filepath.Abs(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	sourceDirAbsPath := path.Dir(sourceFileAbsPath)
	if len(flags.RootDir) != 0 {
		sourceDirAbsPath = flags.RootDir
	}
	controlData, err := getControlData(project, sourceDirAbsPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i := range controlData.subfiles {
		err = project.config.Append(controlData.subfiles[i])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
	project.config.DeleteSection(controlSectionName)

	err = validateProject(project)
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
	controlData.fields[profilesKey] = flags.Profiles
	controlData.fields[overridesKey] = flags.Overrides
	controlData.fields[buildConfigKey] = []string{sourceFileAbsPath}

	cacheDir, cacheFile := getCachePath(flags)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	workDirAbsPath, err := filepath.Abs(flags.getBuildDir())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	controlData.fields[sourceDirKey] = []string{sourceDirAbsPath}
	controlData.fields[workDirKey] = []string{workDirAbsPath}

	components := project.Components()
	for i := range components {
		component, found := project.GetComponent(components[i])
		if !found {
			log.Fatalf("Error: %v", errCouldntLoadComponent)
		}
		component.SetValue(sourceDirKey, []string{path.Join(sourceDirAbsPath, component.Name())})
		component.SetValue(workDirKey, []string{path.Join(workDirAbsPath, component.Name())})
	}
	err = applyControlData(project, controlData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	outConfig, err := os.Create(cacheFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer outConfig.Close()
	err = project.write(outConfig)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = dpl.WriteProject(cacheDir, loaderString)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type ReconfigureFlags struct {
	Append    bool
	Overrides []string
	Profiles  []string
}

func DoReconfigure(flags ReconfigureFlags, args []string) {
}

func loadExistingProject(cacheDir string) (dpl.Project, error) {
	configPath := path.Join(cacheDir, cacheFileName)
	project, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	controlComponent, found := project.getAnyComponent(controlSectionName)
	if !found {
		return nil, errNoControlData
	}
	srcDir := controlComponent.GetValue(sourceDirKey)
	if len(srcDir) != 1 {
		return nil, errMissingSourceDir
	}
	project.srcDir = srcDir[0]
	workDir := controlComponent.GetValue(workDirKey)
	if len(workDir) != 1 {
		return nil, errMissingWorkDir
	}
	project.workDir = workDir[0]
	return project, err
}

func init() {
	err := dpl.RegisterLoader(loaderString, loadExistingProject)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
