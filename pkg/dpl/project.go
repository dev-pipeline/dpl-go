package dpl

import (
	"fmt"
	"os"
	"path"
)

type Component interface {
	Name() string
	ValueNames() []string
	GetValue(string) []string
	ExpandValue(string) ([]string, error)
	SetValue(string, []string)
	EraseValue(string)
	GetSourceDir() string
	GetWorkDir() string
}

type Project interface {
	GetComponent(string) (Component, bool)
	Components() []string
	Write() error
}

type ProjectLoader func(string) (Project, error)

var (
	errCantFindDplDir    error                    = fmt.Errorf("can't find dev-pipeline control directory")
	errAlreadyRegistered error                    = fmt.Errorf("loader identifier string is already registered")
	availableManagers    map[string]ProjectLoader = map[string]ProjectLoader{}
)

const (
	dplDir                 string = ".dpl"
	projectManagerFilename string = "manager"
)

func findProject(startDir string) (string, error) {
	for {
		testDir := path.Join(startDir, dplDir)
		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			return startDir, nil
		}
		nextDir := path.Dir(startDir)
		if startDir == nextDir {
			return "", errCantFindDplDir
		}
		startDir = nextDir
	}
}

func getProjectLoader(dplPath string, managers map[string]ProjectLoader) (ProjectLoader, error) {
	managerPath := path.Join(dplPath, projectManagerFilename)
	data, err := os.ReadFile(managerPath)
	if err != nil {
		return nil, err
	}
	if loader, ok := managers[string(data)]; ok {
		return loader, nil
	}
	return nil, fmt.Errorf("unknown project loader '%v'", string(data))
}

func LoadProject() (Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rootDir, err := findProject(cwd)
	controlDir := path.Join(rootDir, dplDir)
	if err != nil {
		return nil, err
	}
	loader, err := getProjectLoader(controlDir, availableManagers)
	if err != nil {
		return nil, err
	}
	return loader(controlDir)
}

func WriteProject(controlDir string, loaderIdentifier string) error {
	controlFilePath := path.Join(controlDir, projectManagerFilename)
	loaderFile, err := os.Create(controlFilePath)
	if err != nil {
		return err
	}
	defer loaderFile.Close()

	totalWritten := 0
	for {
		count, err := loaderFile.WriteString(loaderIdentifier[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += count
		if totalWritten == len(loaderIdentifier) {
			return nil
		}
	}
}

func RegisterLoader(identifier string, loader ProjectLoader) error {
	_, ok := availableManagers[identifier]
	if ok {
		return errAlreadyRegistered
	}
	availableManagers[identifier] = loader
	return nil
}
