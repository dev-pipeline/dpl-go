package build

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type noneBuilder struct {
}

func (nb noneBuilder) Configure() error {
	return nil
}

func (nb noneBuilder) Build() error {
	return nil
}

func (nb noneBuilder) Install() error {
	return nil
}

func init() {
	RegisterBuilder("none", func(dpl.Component) (Builder, error) {
		return &noneBuilder{}, nil
	})
}
