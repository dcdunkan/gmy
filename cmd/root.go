package cmd

import (
	"errors"
	"fmt"
	"gmy/internal"
	"io/ioutil"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gmy [project name]",
	Short: "Create grammY projects from the command line!",
	Long: `grammY is an amazing Telegram Bot framework for Deno and Node.js.
This application is a tool to generate the needed files to quickly
create a Telegram bot powered by grammY.

Documentation: https://grammy.dev/
GitHub: https://github.com/grammyjs/grammY
Telegram: https://telegram.me/grammyjs`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var project string

		if len(args) > 0 {
			err := projectNameValidator(args[0])
			if err == nil {
				project = args[0]
			}
		}

		if project == "" {
			projectNamePrompt := textinput.New(" > Enter a project name:")
			projectNamePrompt.InitialValue = "my-bot"
			projectNamePrompt.Placeholder = "Project name cannot be empty"
			projectNamePrompt.Validate = projectNameValidator
			input, err := projectNamePrompt.RunPrompt()
			project = input
			if err != nil {
				os.Exit(1)
			}
		}

		platformPrompt := selection.New(" > Choose a platform:",
			selection.Choices([]string{"Deno", "Node", "Other"}))
		platformPrompt.FilterPrompt = "   Select"
		platformPrompt.FilterPlaceholder = "Filter platforms"
		platformPrompt.PageSize = 5

		choice, err := platformPrompt.RunPrompt()
		if err != nil {
			os.Exit(1)
		}
		platform := choice.String

		fmt.Printf(" : Getting templates for %s...\n", platform)

		internal.SetupTemplates(platform, project)

		fmt.Println(`
 + Project folder set up successfully!
   Here are some other configuration that might help you:
 `)

		if platform == "Node" {
			internal.SetupNode(project)
		}

		// Git init
		gitPrompt := confirmation.New(" ? Initialize Git repository", confirmation.Yes)
		initGit, _ := gitPrompt.RunPrompt()

		if initGit {
			internal.SetupGit(project)
		}

		dockerPrompt := confirmation.New(" ? Add Docker related files", confirmation.No)
		docker, _ := dockerPrompt.RunPrompt()

		if docker {
			internal.AddDockerFiles(platform, project)
		}

		pkgInstallPrompt := confirmation.New(" ? Cache/Install dependencies", confirmation.Yes)
		installPkgs, _ := pkgInstallPrompt.RunPrompt()

		if installPkgs {
			internal.InstallDependencies(platform, project)
		}

		fmt.Printf(`
 + Project '%s' created successfully!
 
 Thankyou for using grammY!
 Documentation: https://grammy.dev
 GitHub:        https://github.com/grammyjs/grammY
 Community:     https://telegram.me/grammyjs
`,
			project)
	},
}

func projectNameValidator(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Empty name!")
	}
	if internal.FileExists(name) {
		folder, err := ioutil.ReadDir(name)
		if err != nil {
			return err
		}
		if len(folder) == 0 {
			return nil
		}
		return errors.New("Directory already exists")
	}
	return nil
}

func Execute() {
	rootCmd.DisableSuggestions = false

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
