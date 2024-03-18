package progress

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	_ "embed"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

//go:embed status-page.html.tmpl
var statusPageContent string

func CreateStatusPage(wips []WorkInProgress) ([]byte, error) {
	t := template.New("statusPage")

	t, err := t.Parse(statusPageContent)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, wips); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

const (
	repoURL            = "https://github.com/rhionin/SanderServer.git"
	branchName         = "gh-pages"
	StatusPageFilename = "index.html"
)

func PublishStatusPage(username, apiKey string, content []byte) error {
	// Open the repository (replace with actual URL and authentication)
	auth := &http.BasicAuth{Username: username, Password: apiKey}
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:           repoURL,
		Auth:          auth,
		ReferenceName: branchName,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// Working directory (replace with actual path)
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Open index.html file for writing
	indexFile, err := worktree.Filesystem.Create(StatusPageFilename)
	if err != nil {
		return fmt.Errorf("failed to create/open file: %w", err)
	}

	// Update index.html
	_, err = indexFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Add changes to staging area
	_, err = worktree.Add(StatusPageFilename)
	if err != nil {
		return fmt.Errorf("failed to add file: %w", err)
	}

	// Commit the changes
	commitMsg := "Update " + StatusPageFilename
	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "StormWatch Bot",
			Email: "rhionin@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Push the changes to the remote
	err = repo.Push(&git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}
