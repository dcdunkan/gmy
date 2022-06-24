package internal

import (
	"fmt"
	"os/exec"

	"github.com/erikgeiser/promptkit/selection"
)

func InstallDependencies(platform, dir string) {
	if platform == "Node" {
		pkgManPrompt := selection.New(" > Choose a package manager of your choice:",
			selection.Choices([]string{"npm", "yarn", "pnpm", "None"}))
		pkgManPrompt.FilterPrompt = "   Select"
		pkgManPrompt.FilterPlaceholder = "Filter package managers"
		pkgManPrompt.PageSize = 5

		packageManager, err := pkgManPrompt.RunPrompt()
		if err == nil && packageManager.String != "None" {
			fmt.Println(" : Checking for package manager...")
			cmd := exec.Command(packageManager.String, "--version")
			cmd.Dir = dir
			err := cmd.Run()
			if err != nil {
				fmt.Printf(" ! Skipping... %s not installed.\n", packageManager.String)
			} else {
				fmt.Println(" : Installing dependencies...")
				cmd := exec.Command(packageManager.String, "install")
				cmd.Dir = dir
				err := cmd.Run()

				if err != nil {
					fmt.Printf(" ! Skipping... Failed to install packages.")
				} else {
					fmt.Println(" + Installed dependencies.")
				}
			}
		}
	} else if platform == "Deno" {
		fmt.Println(" : Caching Dependencies...")
		cmd := exec.Command("deno", "cache", "src/mod.ts")
		cmd.Dir = dir
		err := cmd.Run()
		if err != nil {
			fmt.Println(" ! Skipping... Deno CLI not installed.")
		} else {
			fmt.Println(" + Cached dependencies.")
		}
	}
}
