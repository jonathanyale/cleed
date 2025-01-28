package storage

import (
	"os"
	"os/exec"

	"github.com/radulucut/cleed/internal/utils"
)

func (s *LocalStorage) GetExploreRepositoryPath(name string, update bool) (string, error) {
	path, err := s.JoinExploreDir(name)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, cloneRepository(name, path)
		}
		return "", err
	}
	if !update {
		return path, nil
	}
	return path, updateRepository(path)
}

func (s *LocalStorage) RemoveExploreRepository(name string) error {
	path, err := s.JoinExploreDir(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func cloneRepository(url, path string) error {
	out, err := exec.Command("git", "clone", "--depth", "1", "--single-branch", "--no-tags", url, path).CombinedOutput()
	if err != nil {
		return utils.NewInternalError(string(out) + "(" + err.Error() + ")")
	}
	return nil
}

func updateRepository(path string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(currentDir)
	// change the working directory to the repository
	err = os.Chdir(path)
	if err != nil {
		return err
	}
	// fetch the latest changes
	out, err := exec.Command("git", "pull", "--depth", "1", "--no-tags", "origin").CombinedOutput()
	if err != nil {
		return utils.NewInternalError(string(out) + "(" + err.Error() + ")")
	}
	exec.Command("git", "reflog", "expire", "--expire=now", "--all").Run()
	exec.Command("git", "gc", "--prune=now").Run()
	return nil
}
