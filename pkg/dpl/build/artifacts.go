package build

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	errArtifactNotFound error = fmt.Errorf("file not found")
	artifactPattern     *regexp.Regexp
)

const (
	buildArtifactPath string = "build.artifact_path"
	buildInstallPath  string = "build.install_path"
)

func findArtifact(startDir string, filename string) (string, error) {
	entries, err := os.ReadDir(startDir)
	if err != nil {
		return "", err
	}
	for i := range entries {
		if entries[i].Name() == filename {
			return startDir, nil
		}
		if entries[i].IsDir() {
			ret, err := findArtifact(path.Join(startDir, entries[i].Name()), filename)
			if err == nil {
				// found it
				return ret, nil
			}
			if err != errArtifactNotFound {
				// something bad happened
				return "", err
			}
		}
	}
	return "", errArtifactNotFound
}

func findAllArtifacts(component dpl.Component, key string, startDir string) error {
	buildArtifacts, err := component.ExpandValues(key)
	if err != nil {
		return err
	}
	for i := range buildArtifacts {
		groups := artifactPattern.FindStringSubmatch(buildArtifacts[i])
		if groups == nil {
			return errArtifactNotFound
		}
		key := groups[1]
		filename := groups[2]
		fullPath, err := findArtifact(startDir, filename)
		if err != nil {
			return err
		}
		nextKey := fmt.Sprintf("dpl.build.artifact_path.%v", key)
		component.SetValues(nextKey, []string{fullPath})
	}
	return nil
}

func init() {
	var err error
	artifactPattern, err = regexp.Compile(`([a-zA-Z0-9_]+)+=(.+)`)
	if err != nil {
		log.Fatalf("Error compiling pattern: %v", err)
	}
}
