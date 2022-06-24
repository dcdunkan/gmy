package internal

import (
	"fmt"
	"os/exec"
)

func SetupGit(dir string) {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	if err != nil {
		fmt.Println(" - Failed to initialize repository: Git is not installed.")
	} else {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			fmt.Println(" - Failed to initialize a Git repository.")
		} else {
			cmd := exec.Command("git", "add", "-A")
			cmd.Dir = dir
			if err != nil {
				fmt.Println(" - Initialized repository successfully, but failed to add contents to the index.")
			}
		}
	}
}
