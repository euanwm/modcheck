package modcheck

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

// FindGitHub searches deps.dev for a GitHub link.
func FindGitHub(packName string) string {
	uri := "https://deps.dev/_/s/go/p/" + url.QueryEscape(packName) + "/v/"

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return ""
	}

	req = req.WithContext(context.Background())
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`(?m)github.com/[\w-]+/[\w-]+`)

	return re.FindString(string(body))
}
