package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DoCheckGithubActivity(username string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	result, err := getRequest(url)
	if err != nil {
		return err
	}

	countCommits := map[string]int{}
	for _, m := range result {
		var activityName string
		eventType := m["type"].(string)
		repoName := m["repo"].(map[string]interface{})["name"].(string)
		payload := m["payload"].(map[string]interface{})

		switch eventType {
		case "WatchEvent":
			activityName = fmt.Sprintf("Starred %s", repoName)
		case "SponsorshipEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s sponsorship in %s", strings.Title(action), repoName)
		case "ReleaseEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s release in %s", strings.Title(action), repoName)
		case "PushEvent":
			sizeCommits := payload["size"].(float64)
			countCommits[repoName] = countCommits[repoName] + int(sizeCommits)
			continue
		case "PullRequestReviewThreadEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s a pull request review thread in %s", strings.Title(action), repoName)
		case "PullRequestReviewCommentEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s pull request review comment in pull request %s", strings.Title(action), repoName)
		case "PullRequestReviewEvent":
			action := payload["action"].(string)
			pullRequest := payload["pull_request"].(map[string]interface{})
			pullRequestNumber := pullRequest["number"].(float64)
			activityName = fmt.Sprintf("%s pull request review with #%d in %s", strings.Title(action), int(pullRequestNumber), repoName)
		case "PullRequestEvent":
			action := payload["action"].(string)
			number := payload["number"].(float64)
			activityName = fmt.Sprintf("%s pull request #%d in %s", strings.Title(action), int(number), repoName)
		case "PublicEvent":
			activityName = fmt.Sprintf("Made %s public", repoName)
		case "MemberEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s a colaborator in %s", strings.Title(action), repoName)
		case "IssuesEvent":
			action := payload["action"].(string)
			issue := payload["issue"].(map[string]interface{})
			issueNumber := issue["number"].(float64)
			issueTitle := issue["title"].(string)
			activityName = fmt.Sprintf("%s an issue with id #%d in %s with title: %s", strings.Title(action), int(issueNumber), repoName, issueTitle)
		case "IssueCommentEvent":
			action := payload["action"].(string)
			issue := payload["issue"].(map[string]interface{})
			issueNumber := issue["number"].(float64)
			commentID := payload["comment"].(map[string]interface{})["id"].(float64)
			activityName = fmt.Sprintf("%s a new issue comment of issue number #%d with comment id %d in %s", strings.Title(action), int(issueNumber), int(commentID), repoName)
		case "GollumEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s the wiki in %s", strings.Title(action), repoName)
		case "ForkEvent":
			forkee := payload["forkee"].(map[string]interface{})
			activityName = fmt.Sprintf("Forked %s to %s", repoName, forkee["full_name"].(string))
		case "DeleteEvent":
			refType := payload["ref_type"].(string)
			ref := payload["ref"].(string)
			activityName = fmt.Sprintf("Deleted %s %s in %s", refType, ref, repoName)
		case "CreateEvent":
			refType := payload["ref_type"].(string)
			if refType == "repository" {
				activityName = fmt.Sprintf("Created a new %s in %s", refType, repoName)
				continue
			}
			ref := payload["ref"].(string)
			activityName = fmt.Sprintf("Created a new %s %s in %s", refType, ref, repoName)
		case "CommitCommentEvent":
			action := payload["action"].(string)
			activityName = fmt.Sprintf("%s commented on a commit in %s", strings.Title(action), repoName)
		}
		fmt.Println("-", activityName)
	}

	for s, i := range countCommits {
		fmt.Printf("- Pushed %d commits to %s\n", i, s)
	}

	return nil

}

func getRequest(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error fetching data from Github API")
	}

	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("username not found")
		}
		return nil, errors.New("invalid response status from Github API")
	}

	if string(bodyRes) == "[]" {
		return nil, errors.New("no activity found")

	}

	var result []map[string]interface{}
	if err := json.Unmarshal(bodyRes, &result); err != nil {
		return nil, errors.New("no parsing response body")
	}

	return result, nil
}
