package utils

import (
	"errors"
	"os/exec"
)

func InitGit(dir string) error {
	err := IsInstalled("git")

	if err != nil {
		return errors.New("Git is not installed")
	}

	initCmd := exec.Command("git", "init")
	initCmd.Dir = dir
	err = initCmd.Run()

	if err != nil {
		return errors.New("Failed to initialize a Git repository.")
	}

	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = dir
	err = addCmd.Run()

	if err != nil {
		return errors.New("Initialized repository successfully, but failed to add contents to the index.")
	}

	return nil
}
