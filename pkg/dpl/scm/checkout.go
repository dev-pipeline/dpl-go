package scm

import (
	"fmt"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"

	"github.com/dev-pipeline/dpl-go/internal/common"
)

const (
	scmUriKey string = "scm.uri"
)

var (
	CheckoutTask = common.Task{
		Name: "scm",
		Work: checkout,
	}
)

func checkout(component dpl.Component) error {
	scmUris, err := component.ExpandValue(scmUriKey)
	if err != nil {
		return err
	}
	for _, uri := range scmUris {
		scmInfo, err := BuildScmInfo(uri)
		if err != nil {
			return err
		}
		scmBuilder := GetHandler(scmInfo.Scheme)
		if scmBuilder == nil {
			return fmt.Errorf("no handler for %v", scmInfo.Scheme)
		}
		handler, err := scmBuilder(component)
		if err != nil {
			return err
		}
		err = handler.Checkout(scmInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
