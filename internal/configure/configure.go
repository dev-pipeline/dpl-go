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

	loaderString string = "configure"
)

var (
	errMissingSourceDir error = fmt.Errorf("project missing source directory information")
	errMissingWorkDir   error = fmt.Errorf("project missing work directory information")
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

	sourceFileAbsPath, err := filepath.Abs(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	sourceDirAbsPath := path.Dir(sourceFileAbsPath)

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

	workDirAbsPath, err := filepath.Abs(flags.getBuildDir())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defaultComponent.SetValue("dpl.build_config", []string{sourceFileAbsPath})
	defaultComponent.SetValue("dpl.src_dir", []string{sourceDirAbsPath})
	defaultComponent.SetValue("dpl.work_dir", []string{workDirAbsPath})

	outConfig, err := os.Create(cacheFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer outConfig.Close()
	err = project.Write(outConfig)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = dpl.WriteProject(cacheDir, loaderString)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func loadExistingProject(cacheDir string) (dpl.Project, error) {
	configPath := path.Join(cacheDir, cacheFileName)
	project, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	defaultComponent, err := project.getDefaultComponent()
	if err != nil {
		return nil, err
	}
	srcDir := defaultComponent.GetValue("dpl.src_dir")
	if len(srcDir) != 1 {
		return nil, errMissingSourceDir
	}
	project.srcDir = srcDir[0]
	workDir := defaultComponent.GetValue("dpl.work_dir")
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
