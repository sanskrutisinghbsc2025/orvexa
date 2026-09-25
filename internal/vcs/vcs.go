// Package vcs provides a clean wrapper over go-git for Orvexa's
// resume versioning logic.
package vcs

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Repo wraps the go-git repository.
type Repo struct {
	gitRepo *git.Repository
}

// Open opens the Orvexa repository in the current directory.
func Open() (*Repo, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("not a orvexa repo — run 'orvexa init' first")
	}
	return &Repo{gitRepo: r}, nil
}

// Init initializes a new git repository in the current directory.
func Init() (*Repo, error) {
	r, err := git.PlainInit(".", false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}
	return &Repo{gitRepo: r}, nil
}

// CurrentBranch returns the name of the active branch.
func (r *Repo) CurrentBranch() (string, error) {
	head, err := r.gitRepo.Head()
	if err != nil {
		return "(initial branch)", nil
	}
	return head.Name().Short(), nil
}

// ListBranches returns all local branches.
func (r *Repo) ListBranches() ([]string, error) {
	iter, err := r.gitRepo.Branches()
	if err != nil {
		return nil, err
	}
	var names []string
	iter.ForEach(func(ref *plumbing.Reference) error {
		names = append(names, ref.Name().Short())
		return nil
	})
	return names, nil
}

// CreateBranch creates and switches to a new branch.
func (r *Repo) CreateBranch(name string) error {
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
	})
}

// SwitchBranch switches to an existing branch.
func (r *Repo) SwitchBranch(name string) error {
	cmd := exec.Command("git", "switch", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// DeleteBranch deletes a branch by name.
func (r *Repo) DeleteBranch(name string) error {
	if current, _ := r.CurrentBranch(); current == name {
		return fmt.Errorf("cannot delete an active branch")
	}
	return r.gitRepo.Storer.RemoveReference(plumbing.NewBranchReferenceName(name))
}

// Commit saves changes to resume.json with a message.
func (r *Repo) Commit(message string) (string, error) {
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return "", err
	}
	_, err = w.Add("resume.json")
	if err != nil {
		return "", err
	}
	hash, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Orvexa User",
			Email: "user@orvexa.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

// Status returns the checkout status of the repository.
func (r *Repo) Status() (git.Status, error) {
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return nil, err
	}
	return w.Status()
}

// Log returns the commit history iterator.
func (r *Repo) Log() (object.CommitIter, error) {
	return r.gitRepo.Log(&git.LogOptions{})
}

// Restore reverts resume.json to a specific commit hash.
func (r *Repo) Restore(shortHash string, force bool) (string, error) {
	fullHash, err := r.gitRepo.ResolveRevision(plumbing.Revision(shortHash))
	if err != nil {
		return "", fmt.Errorf("version [%s] not found", shortHash)
	}

	// Only protect against uncommitted changes to resume.json.
	// Use system Git here because it gives reliable worktree status
	// on Windows/Git Bash.
	statusCmd := exec.Command("git", "status", "--porcelain", "--", "resume.json")
	output, err := statusCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to check resume status: %s", string(output))
	}

	if !force && len(output) > 0 {
		return "", fmt.Errorf("unsaved changes detected; use --force to overwrite")
	}

	// Restore only resume.json from the requested version.
	// This keeps the user on the current branch.
	checkoutCmd := exec.Command(
		"git",
		"checkout",
		fullHash.String(),
		"--",
		"resume.json",
	)

	output, err = checkoutCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to restore version: %s", string(output))
	}

	return fullHash.String(), nil
}

// Sync performs a git rebase of the current branch onto the target branch.
// Currently wraps system git due to complexity of rebase in go-git.
func (r *Repo) Sync(targetBranch string) ([]byte, error) {
	return exec.Command("git", "rebase", targetBranch).CombinedOutput()
}
