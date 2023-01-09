package scm

import (
	"fmt"
	"net/url"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ScmInfo struct {
	Scheme    string
	Path      string
	Arguments map[string]string
}

type ScmHandler interface {
	Checkout(ScmInfo) error
}

type MakeScm func(dpl.Component) (ScmHandler, error)

var (
	scms map[string]MakeScm = map[string]MakeScm{}
)

func AddHandler(protocol string, handler MakeScm) error {
	scms[protocol] = handler
	return nil
}

func GetHandler(protocol string) MakeScm {
	handler, found := scms[protocol]
	if found {
		return handler
	}
	return nil
}

func BuildScmInfo(uriInfo string) (ScmInfo, error) {
	scmInfo := ScmInfo{}
	uri, err := url.Parse(uriInfo)
	if err != nil {
		return scmInfo, err
	}

	trimmedPath, arguments, err := extractArguments(uri.Path, ";")
	if err != nil {
		return scmInfo, err
	}

	userPrefix := ""
	if uri.User != nil {
		userPrefix = fmt.Sprintf("%v@", uri.User)
	}
	scmInfo.Scheme = uri.Scheme
	scmInfo.Path = fmt.Sprintf("%v%v%v", userPrefix, uri.Host, trimmedPath)
	scmInfo.Arguments = arguments

	return scmInfo, nil
}
