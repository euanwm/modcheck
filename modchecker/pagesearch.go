package modchecker

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
)

func FindGitHub(packName string) (gh string) {
	urlStr := "https://deps.dev/_/s/go/p/" + url.QueryEscape(packName) + "/v/"
	resp, err := http.Get(urlStr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`(?m)github.com/[\w-]+/[\w-]+`)
	gh = re.FindString(string(body))
	return gh
}
