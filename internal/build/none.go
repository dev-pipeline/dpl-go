package build

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type noneBuilder struct {
}

func (nb noneBuilder) Configure(*BuildConfig) error {
	return nil
}

func (nb noneBuilder) Build(*BuildConfig) error {
	return nil
}

func (nb noneBuilder) Install(destdir string) error {
	return nil
}

func init() {
	RegisterBuilder("none", func(dpl.Component) (Builder, error) {
		return &noneBuilder{}, nil
	})
}
