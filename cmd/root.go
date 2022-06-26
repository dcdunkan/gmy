package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dcdunkan/gmy/internal/files"
	t "github.com/dcdunkan/gmy/internal/templates"
	"github.com/dcdunkan/gmy/internal/utils"

	"github.com/cavaliergopher/grab/v3"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/spf13/cobra"
)

var footer = `
 + Project '%s' created successfully!
 
 Thankyou for using grammY <3
 . Documentation: https://grammy.dev
 . GitHub:        https://github.com/grammyjs/grammY
 . Community:     https://telegram.me/grammyjs
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gmy [project name]",
	Short: "Create grammY projects from the command line!",
	Long: `grammY is an amazing Telegram Bot framework for Deno and Node.js.
This application is a tool to generate the needed files to quickly
create a Telegram bot powered by grammY.

https://grammy.dev/`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string

		if len(args) > 0 {
			err := ValidateProjectName(args[0])
			if err == nil {
				projectName = args[0]
			}
		}

		if projectName == "" {
			prompt := textinput.New(" > Enter a name for your project:")
			prompt.Placeholder = "Project name cannot be empty"
			prompt.Validate = ValidateProjectName

			input, err := prompt.RunPrompt()
			if err != nil {
				os.Exit(1)
			}

			projectName = input
		}

		platformPrompt := selection.New(
			" > Choose a platform:",
			selection.Choices([]string{"Deno", "Node", "Other"}),
		)
		platformPrompt.FilterPrompt = "   Platforms"
		platformPrompt.FilterPlaceholder = "Type to filter"
		platformPrompt.PageSize = 5

		choice, err := platformPrompt.RunPrompt()

		if err != nil {
			os.Exit(1)
		}

		platform := choice.String

		fmt.Printf(" : Getting templates for %s...\n", platform)
		templates, err := t.GetTemplates(platform)

		if err != nil {
			fmt.Println(" ! Failed to fetch templates :(")
			os.Exit(1)
		}

		var templateChoices []string

		for i := 0; i < len(templates); i++ {
			templateChoices = append(templateChoices, templates[i].Name)
		}

		templatePrompt := selection.New(" > Choose a template:", selection.Choices(templateChoices))
		templatePrompt.FilterPrompt = "   Templates for " + platform
		templatePrompt.FilterPlaceholder = "Type to filter"
		templatePrompt.PageSize = 5

		choice, err = templatePrompt.RunPrompt()

		if err != nil {
			os.Exit(1)
		}

		var selectedTemplate t.Template

		for _, template := range templates {
			if template.Name == choice.String {
				selectedTemplate = template
				break
			}
		}

		fmt.Printf(
			" : Downloading template assets from %s/%s...\n",
			selectedTemplate.Owner, selectedTemplate.Repository,
		)

		tempDir, err := os.MkdirTemp("", "gmy-template-")
		defer os.RemoveAll(tempDir)

		downloadUrl, err := t.GetDownloadUrl(selectedTemplate)

		if err != nil {
			fmt.Println(" ! Failed to fetch template details :(")
			os.Exit(1)
		}

		resp, _ := grab.Get(tempDir, downloadUrl)

		if err := resp.Err(); err != nil {
			fmt.Printf(" ! Failed to download assets from %s\n", downloadUrl)
			os.Exit(1)
		}

		fmt.Println(" + Assets downloaded")
		fmt.Println(" : Setting up project folder...")

		file, err := os.Open(resp.Filename)
		err = files.Untar(projectName, file)

		if err != nil {
			fmt.Println(" ! Failed to extract downloaded assets: ", err)
			os.Exit(1)
		}

		fmt.Printf(`
 + Project folder set up successfully!
   Additional configuration options for you:

`)

		gitPrompt := confirmation.New(" ? Initialize Git repository", confirmation.Yes)
		initGit, _ := gitPrompt.RunPrompt()

		if initGit {
			err := utils.InitGit(projectName)

			if err != nil {
				fmt.Println(" - Failed to initialize Git repository.")
			} else {
				fmt.Println(" + Git repository initialized.")
			}
		}

		if selectedTemplate.TSConfig_prompt {
			tscPrompt := confirmation.New(" ? Would you like to add a tsconfig.json", confirmation.No)
			tsc, _ := tscPrompt.RunPrompt()
			if tsc {
				_, err := grab.Get(projectName, "https://raw.githubusercontent.com/dcdunkan/gmy/main/configs/tsconfig.json")
				if err != nil {
					fmt.Println(" - Skipping... Failed to add a tsconfig.json.")
				}
			}
		}

		if selectedTemplate.Docker_prompt {
			dockerPrompt := confirmation.New(" ? Add Docker related files", confirmation.No)
			docker, err := dockerPrompt.RunPrompt()

			if err == nil && docker {
				err := files.AddDockerFiles(platform, projectName)
				if err != nil {
					fmt.Println(" - Failed to add Docker files.")
				}
			}
		}

		if platform == "Other" {
			fmt.Printf(footer, projectName)
			return
		}

		pkgInstallPrompt := confirmation.New(" ? Cache/Install dependencies", confirmation.Yes)
		installPkgs, err := pkgInstallPrompt.RunPrompt()

		if err != nil {
			os.Exit(1)
		}

		if installPkgs {
			if platform == "Deno" {
				fmt.Println(" : Caching dependencies...")
				err := files.CacheDenoDeps(projectName, selectedTemplate.Cache_file)

				if err != nil {

					fmt.Println(" ! Failed to cache dependencies.")
				} else {
					fmt.Println(" + Cached the dependencies.")
				}
			}

			if platform == "Node" {
				prompt := selection.New(
					" > Choose a package manager of your choice:",
					selection.Choices([]string{"npm", "yarn", "pnpm", "None"}),
				)

				prompt.FilterPrompt = "   Package Managers"
				prompt.FilterPlaceholder = "Type to filter"
				prompt.PageSize = 5

				choice, err := prompt.RunPrompt()
				if err == nil && choice.String != "None" {
					packageManager := choice.String
					fmt.Println(" : Checking for package manager...")
					err := utils.IsInstalled(packageManager)

					if err != nil {
						fmt.Println(" ! Seems like the package manager is not installed.")
					} else {
						fmt.Println(" : Installing dependencies...")
						err := files.InstallNodeDeps(packageManager, projectName)
						if err != nil {
							fmt.Println(" ! Failed to install dependencies.")
						} else {
							fmt.Println(" + Dependencies installed.")
						}
					}
				}

				fmt.Printf(" : Updating package name...\n")
				files.UpdatePackageName(projectName)
			}
		}

		fmt.Printf(footer, projectName)
	},
}

func ValidateProjectName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Empty name!")
	}

	if files.Exists(name) {
		folder, err := os.ReadDir(name)

		if err != nil {
			return err
		}

		if len(folder) == 0 {
			return nil
		}

		return errors.New("Directory already exists!")
	}

	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
