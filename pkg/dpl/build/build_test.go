package build

import (
	"fmt"
	"log"
	"testing"

	"github.com/dev-pipeline/dpl-go/internal/test/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func TestDoBuild(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {"none"},
		},
	}
	err := doFullBuild(c)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestNoBuilder(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{},
	}
	err := doFullBuild(c)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestBuildConfigureError(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {configureErrorBuilder},
		},
	}
	err := doFullBuild(c)
	if err != errConfigureError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestBuildBuildError(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {buildErrorBuilder},
		},
	}
	err := doFullBuild(c)
	if err != errBuildError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestBuildInstallError(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {installErrorBuilder},
		},
	}
	err := doFullBuild(c)
	if err != errInstallError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestBuildInstallErrorNoInstall(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey:     {installErrorBuilder},
			installMethodKey: {"none"},
		},
	}
	err := doFullBuild(c)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

var (
	errConfigureError error = fmt.Errorf("configure error")
	errBuildError     error = fmt.Errorf("build error")
	errInstallError   error = fmt.Errorf("install error")
)

const (
	configureErrorBuilder string = "configure-error-builder"
	buildErrorBuilder     string = "build-error-builder"
	installErrorBuilder   string = "install-error-builder"
)

type errorBuilder struct {
	configureErr error
	buildErr     error
	installErr   error
}

func (eb errorBuilder) Configure(*BuildConfig) error {
	return eb.configureErr
}

func (eb errorBuilder) Build(*BuildConfig) error {
	return eb.buildErr
}

func (eb errorBuilder) Install(string) error {
	return eb.installErr
}

func makeConfigureErrorBuilder(dpl.Component) (Builder, error) {
	return &errorBuilder{
		configureErr: errConfigureError,
	}, nil
}

func makeBuildErrorBuilder(dpl.Component) (Builder, error) {
	return &errorBuilder{
		buildErr: errBuildError,
	}, nil
}

func makeInstallErrorBuilder(dpl.Component) (Builder, error) {
	return &errorBuilder{
		installErr: errInstallError,
	}, nil
}

func init() {
	err := RegisterBuilder(configureErrorBuilder, makeConfigureErrorBuilder)
	if err != nil {
		log.Fatalf("Error registring builder: %v", err)
	}
	err = RegisterBuilder(buildErrorBuilder, makeBuildErrorBuilder)
	if err != nil {
		log.Fatalf("Error registring builder: %v", err)
	}
	err = RegisterBuilder(installErrorBuilder, makeInstallErrorBuilder)
	if err != nil {
		log.Fatalf("Error registring builder: %v", err)
	}
}
