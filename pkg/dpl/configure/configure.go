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

	customSourceDirKey string = "source_dir"
)

var (
	errNoControlData    error = fmt.Errorf("no control data")
	errWrongProjectType error = fmt.Errorf("project wasn't created using 'configure'")
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

func configureFromScratch(flags ConfigureFlags) (*IniProject, error) {
	project, err := loadConfig(flags.ConfigFile)
	if err != nil {
		return nil, err
	}

	sourceFileAbsPath, err := filepath.Abs(flags.ConfigFile)
	if err != nil {
		return nil, err
	}
	sourceDirAbsPath := path.Dir(sourceFileAbsPath)
	controlData, err := getControlData(project, sourceDirAbsPath)
	if len(flags.RootDir) != 0 {
		sourceDirAbsPath = flags.RootDir
	}
	if err != nil {
		return nil, err
	}

	for i := range controlData.subfiles {
		err = project.config.Append(controlData.subfiles[i])
		if err != nil {
			return nil, err
		}
	}
	project.config.DeleteSection(controlSectionName)

	err = validateProject(project)
	if err != nil {
		return nil, err
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dplConfigDir := path.Join(homedir, ".dev-pipeline")

	modSet, err := loadProfiles(dplConfigDir, flags.Profiles)
	if err != nil {
		return nil, err
	}
	err = applyProfiles(modSet, project)
	if err != nil {
		return nil, err
	}
	err = applyOverrides(dplConfigDir, flags.Profiles, project)
	if err != nil {
		return nil, err
	}
	controlData.fields[profilesKey] = flags.Profiles
	controlData.fields[overridesKey] = flags.Overrides
	controlData.fields[buildConfigKey] = []string{sourceFileAbsPath}

	cacheDir, cacheFile := getCachePath(flags)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, err
	}

	workDirAbsPath, err := filepath.Abs(flags.getBuildDir())
	if err != nil {
		return nil, err
	}
	controlData.fields[sourceDirKey] = []string{sourceDirAbsPath}
	controlData.fields[workDirKey] = []string{workDirAbsPath}

	components := project.ComponentNames()
	for i := range components {
		component, err := project.getConfigComponent(components[i])
		if err != nil {
			return nil, err
		}
		sourceDir, err := component.GetSingleValue(customSourceDirKey)
		if err != nil {
			if err != errMissingKey {
				return nil, err
			}
			// error is just a missing key, so nothing to worry about; use default source directory
			sourceDir = path.Join(sourceDirAbsPath, component.Name())
		} else {
			if !path.IsAbs(sourceDir) {
				sourceDir = path.Join(sourceDirAbsPath, sourceDir)
			}
		}
		component.SetValues(sourceDirKey, []string{sourceDir})
		component.SetValues(workDirKey, []string{path.Join(workDirAbsPath, component.Name())})
	}
	err = applyControlData(project, controlData)
	if err != nil {
		return nil, err
	}

	outConfig, err := os.Create(cacheFile)
	if err != nil {
		return nil, err
	}
	defer outConfig.Close()
	err = project.write(outConfig)
	if err != nil {
		return nil, err
	}
	err = dpl.WriteProject(cacheDir, loaderString)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func DoConfigure(flags ConfigureFlags, args []string) error {
	_, err := configureFromScratch(flags)
	return err
}

type ReconfigureFlags struct {
	Append    bool
	Overrides []string
	Profiles  []string
}

func DoReconfigure(flags ReconfigureFlags, args []string) error {
	project, err := dpl.LoadProject()
	if err != nil {
		return err
	}
	realProject, ok := project.(*IniProject)
	if !ok {
		return errWrongProjectType
	}
	controlComponent, err := realProject.getAnyComponent(controlSectionName)
	if err != nil {
		return errNoControlData
	}
	cf := ConfigureFlags{}
	err = setupFlags(&cf, controlComponent)
	if err != nil {
		return err
	}

	if len(flags.Overrides) != 0 {
		if flags.Append {
			cf.Overrides = append(cf.Overrides, flags.Overrides...)
		} else {
			cf.Overrides = flags.Overrides
		}
	}
	if len(flags.Profiles) != 0 {
		if flags.Append {
			cf.Profiles = append(cf.Profiles, flags.Profiles...)
		} else {
			cf.Profiles = flags.Profiles
		}
	}
	project, err = configureFromScratch(cf)
	return err
}

func loadExistingProject(cacheDir string) (dpl.Project, error) {
	configPath := path.Join(cacheDir, cacheFileName)
	project, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}

	controlComponent, err := project.getAnyComponent(controlSectionName)
	if err != nil {
		return nil, errNoControlData
	}
	config, err := projectUpToDate(info.ModTime(), controlComponent)
	if err != nil {
		if err != errProjectOutOfDate {
			return nil, err
		}
		project, err = configureFromScratch(config)
		if err != nil {
			return nil, err
		}
		project.configFile = configPath
	}
	return project, err
}

func init() {
	err := dpl.RegisterLoader(loaderString, loadExistingProject)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
