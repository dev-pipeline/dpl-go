package git

import (
	"path"

	gogit "github.com/go-git/go-git/v5"

	"github.com/dev-pipeline/dpl-go/pkg/dpl/scm"
)

func gitClone(srcDir string, info scm.ScmInfo) (*gogit.Repository, error) {
	r, err := gogit.PlainClone(srcDir, false, &gogit.CloneOptions{
		URL:        info.Path,
		NoCheckout: true,
	})
	return r, err
}

func getRepository(srcDir string, info scm.ScmInfo) (*gogit.Repository, error) {
	if offset, found := info.Arguments["path"]; found {
		srcDir = path.Join(srcDir, offset)
	}
	r, err := gogit.PlainOpen(srcDir)
	if err == gogit.ErrRepositoryNotExists {
		return gitClone(srcDir, info)
	}
	if err != nil {
		// some other error
		return nil, err
	}

	// we've got a repo, so update it
	err = r.Fetch(&gogit.FetchOptions{})
	if err != nil && err != gogit.NoErrAlreadyUpToDate {
		return nil, err
	}
	return r, nil
}
