package service

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hulizhen/blogo/config"
)

type RepoService struct {
	config *config.Config
}

func NewRepoService(cfg *config.Config) *RepoService {
	return &RepoService{config: cfg}
}

func (r *RepoService) UpdateRepo() error {
	localPath := r.config.Repository.LocalPath
	repo := r.config.Repository

	// Clone the blog repo if it doesn't exist.
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		cmd := fmt.Sprintf("git clone %v %v", repo.RemoteURL, repo.LocalPath)
		err = exec.Command("/bin/sh", "-c", cmd).Run()
		if err != nil {
			return err
		}
	}

	// Check out branch and pull latest changes of the blog repo.
	cmd := fmt.Sprintf("cd %v && git checkout %v && git pull", repo.LocalPath, repo.Branch)
	err := exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		return err
	}

	// Checkout to the target branch.
	cmd = fmt.Sprintf("cd %v && git checkout %v", repo.LocalPath, repo.Branch)
	err = exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		return err
	}

	return nil
}
