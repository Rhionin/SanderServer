package progress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	_ "embed"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

	useMemfs := false

	cloneOpts := &git.CloneOptions{
		URL:           repoURL,
		Auth:          auth,
		ReferenceName: branchName,
	}

	var repo *git.Repository
	var err error
	if useMemfs {
		repo, err = git.Clone(memory.NewStorage(), memfs.New(), cloneOpts)
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
	} else {
		tmpDir := "/Users/cjc/go/src/github.com/Rhionin/SanderServer/tmp"
		err = os.RemoveAll(tmpDir)
		if err != nil {
			return fmt.Errorf("failed to remove tmp dir: %w", err)
		}
		repo, err = git.PlainClone(tmpDir, false, cloneOpts)
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
	}

	// Working directory (replace with actual path)
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}
	fmt.Println("head before change", headRef.Hash(), headRef.Name())
	PrintDiffForCommit(repo, headRef.Hash())

	// fmt.Println("worktree.Filesystem.Root()", worktree.Filesystem.Root())
	// files, err := worktree.Filesystem.ReadDir(worktree.Filesystem.Root())
	// if err != nil {
	// 	return err
	// }

	// for _, d := range files {
	// 	fmt.Println("\t", d)
	// }

	// // Debug: Print original file contents
	// indexFile, err := worktree.Filesystem.Open(StatusPageFilename)
	// if err != nil {
	// 	return fmt.Errorf("failed to open original index.html file: %w", err)
	// }
	// originalFileContents, err := io.ReadAll(indexFile)
	// if err != nil {
	// 	return fmt.Errorf("failed to read original index.html file: %w", err)
	// }
	// fmt.Println("Original file contents:\n", string(originalFileContents))

	// Open index.html file for writing
	indexFile, err := worktree.Filesystem.Create(StatusPageFilename)
	if err != nil {
		return fmt.Errorf("failed to create/open file: %w", err)
	}

	fmt.Println("Updated page content:", string(content))
	// fmt.Println("Updated page content:", string(content))

	// Update index.html
	_, err = indexFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Debug: Print original file contents
	f, err := worktree.Filesystem.Open(StatusPageFilename)
	if err != nil {
		return fmt.Errorf("failed to open original index.html file: %w", err)
	}
	originalFileContents, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read original index.html file: %w", err)
	}
	fmt.Println("updated file contents:\n", string(originalFileContents))

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

	headRef, err = repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}
	fmt.Println("head after change", headRef.Hash(), headRef.Name())
	PrintDiffForCommit(repo, headRef.Hash())

	// Push the changes to the remote
	err = repo.Push(&git.PushOptions{
		Auth: auth,
		// RemoteName: "origin",
		// RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:%s", branchName, branchName))},
		// Force:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

func PrintDiffForCommit(repo *git.Repository, commitHash plumbing.Hash) error {
	// Get the commit object
	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		return fmt.Errorf("failed to get commit object: %w", err)
	}
	fmt.Println("ASDF", commit)

	// Get the tree objects for parent and current commit
	tree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get commit tree: %w", err)
	}

	// If there's a single parent, get the parent tree
	var parentTree *object.Tree
	if commit.NumParents() > 0 {
		parents := commit.Parents()
		parentCommit, err := parents.Next()
		if err != nil {
			return fmt.Errorf("failed to get parent commit object: %w", err)
		}
		parentTree, err = parentCommit.Tree()
		if err != nil {
			return fmt.Errorf("failed to get parent commit tree: %w", err)
		}
	}

	// Use go-diff/diffmatchpatch for diff generation
	diff, err := tree.Diff(parentTree)
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	// Print the diff to standard output
	for _, patch := range diff {
		fmt.Println(patch)
	}

	return nil
}
