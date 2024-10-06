package main

import (
	"fmt"
	"github-activity-cli/github"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a Github username")
		return
	}

	username := os.Args[1]
	fmt.Println("GitHub User Activity of", username)
	if err := github.DoCheckGithubActivity(username); err != nil {
		fmt.Println(err)
	}
}
