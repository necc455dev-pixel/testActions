package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

// pullSummary keeps the API response minimal while exposing useful fields.
type pullSummary struct {
	Number  int       `json:"number"`
	Title   string    `json:"title"`
	State   string    `json:"state"`
	URL     string    `json:"url"`
	User    string    `json:"user"`
	Updated time.Time `json:"updated_at"`
}

func main() {
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")
	token := os.Getenv("GITHUB_TOKEN")

	if owner == "" || repo == "" {
		log.Fatal("GITHUB_OWNER and GITHUB_REPO environment variables are required")
	}

	client := newGitHubClient(token)

	http.HandleFunc("/pulls", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		prs, err := listPullRequests(ctx, client, owner, repo)
		if err != nil {
			log.Printf("failed to list pull requests: %v", err)
			http.Error(w, "failed to list pull requests", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(prs); err != nil {
			log.Printf("failed to write response: %v", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	addr := ":8080"
	log.Printf("Serving pull requests for %s/%s on %s", owner, repo, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func newGitHubClient(token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)
	return github.NewClient(httpClient)
}

func listPullRequests(ctx context.Context, client *github.Client, owner, repo string) ([]pullSummary, error) {
	opts := &github.PullRequestListOptions{State: "open", Sort: "updated", Direction: "desc"}
	pulls, _, err := client.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	result := make([]pullSummary, 0, len(pulls))
	for _, pr := range pulls {
		result = append(result, pullSummary{
			Number:  pr.GetNumber(),
			Title:   pr.GetTitle(),
			State:   pr.GetState(),
			URL:     pr.GetHTMLURL(),
			User:    pr.GetUser().GetLogin(),
			Updated: pr.GetUpdatedAt().Time,
		})
	}
	return result, nil
}
