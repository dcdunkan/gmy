package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/erikgeiser/promptkit/selection"
)

type Template struct {
	Name, Git_url, Type, Download_url string
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func SetupTemplates(platform, project string) {
	templates := []Template{}
	url := fmt.Sprintf("https://api.github.com/repos/grammyjs/create-grammy/contents/templates/%s", strings.ToLower(platform))
	jsonErr := getJson(url, &templates)

	if jsonErr != nil {
		fmt.Println(" ! Failed to fetch templates :(")
		os.Exit(1)
	}

	choices := []string{}
	for i := 0; i < len(templates); i++ {
		if templates[i].Name[0:1] != "." {
			choices = append(choices, templates[i].Name)
		}
	}
	templateName := selection.New(" > Choose a template:", selection.Choices(choices))
	templateName.FilterPrompt = "   Select"
	templateName.FilterPlaceholder = "Filter templates"
	template, err := templateName.RunPrompt()
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf(" : Downloading template assets from %s...\n", template.String)
	tempDir, err := os.MkdirTemp("", "gmy-template-")
	defer os.RemoveAll(tempDir)

	var toFetch string

	for i := 0; i < len(templates); i++ {
		if template.String != templates[i].Name {
			continue
		}
		url := strings.Replace(templates[i].Git_url, "https://api.github.com/repos/", "", 1)
		split := strings.Split(url, "/")
		repository := strings.Join(split[0:2], "/")
		treeId := split[len(split)-1]
		toFetch = fmt.Sprintf("https://codeload.github.com/%s/tar.gz/%s", repository, treeId)
		break
	}

	resp, _ := grab.Get(tempDir, toFetch)

	if err := resp.Err(); err != nil {
		fmt.Printf(" ! Failed to download assets from \n%s\n", url)
		os.Exit(1)
	}

	fmt.Println(" + Assets downloaded")
	fmt.Println(" : Setting up project folder...")

	f, err := os.Open(resp.Filename)
	err = Untar(project, f)
	if err != nil {
		fmt.Println(" ! Failed to extract downloaded assets: ", err)
		os.Exit(1)
	}
}
