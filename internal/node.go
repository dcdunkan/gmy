package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/erikgeiser/promptkit/confirmation"
)

func SetupNode(project string) {
	projectFolder, err := ioutil.ReadDir(project)
	if err != nil {
		fmt.Println(" ! Failed to read project folder.")
		os.Exit(1)
	}

	var tscExists bool
	var packageJsonExists bool
	for i := 0; i < len(projectFolder); i++ {
		if projectFolder[i].Name() == "package.json" {
			packageJsonExists = true
		}
		if len(projectFolder[i].Name()) > 9 && projectFolder[i].Name()[0:9] == "tsconfig." {
			tscExists = true
		}
	}

	// Add tsconfig.json
	if !tscExists {
		tscPrompt := confirmation.New(" ? Would you like to add a tsconfig.json", confirmation.No)
		tsc, _ := tscPrompt.RunPrompt()
		if tsc {
			_, err := grab.Get(project, "https://raw.githubusercontent.com/grammyjs/create-grammy/main/configs/tsconfig.json")
			if err != nil {
				fmt.Println(" - Skipping... Failed to add a tsconfig.json.")
			}
		}
	}

	// Update package name in package.json
	fmt.Println(" : Updating package name in package.json...")
	if packageJsonExists {
		pkgJsonPath := filepath.Join(project, "package.json")
		bytes, err := os.ReadFile(pkgJsonPath)
		if err == nil {
			packageJson := make(map[string]interface{})
			json.Unmarshal(bytes, &packageJson)
			packageJson["name"] = project
			data, _ := json.Marshal(packageJson)
			err := os.WriteFile(pkgJsonPath, data, 0777)
			if err != nil {
				fmt.Println(" - Skipping... Failed to write package.json.")
			}
		} else {
			fmt.Println(" - Skipping... Failed to write package.json.")
		}
	}
}
