package modcheck

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

// Project is a struct that contains information about a project.
type Project struct {
	ProjectKey struct {
		ID string `json:"id"`
	} `json:"projectKey"`

	OpenIssuesCount string `json:"openIssuesCount"`
	StarsCount      string `json:"starsCount"`
	ForksCount      string `json:"forksCount"`
	License         string `json:"license"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	Scorecard       struct {
		Date       string `json:"date"`
		Repository struct {
			Name   string `json:"name"`
			Commit string `json:"commit"`
		} `json:"repository"`
		Scorecard struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
		} `json:"scorecard"`
		Checks []struct {
			Name          string `json:"name"`
			Documentation struct {
				ShortDescription string `json:"shortDescription"`
				URL              string `json:"url"`
			} `json:"documentation"`
			Score   string   `json:"score"`
			Reason  string   `json:"reason"`
			Details []string `json:"details"`
		} `json:"checks"`
		OverallScore float32  `json:"overallScore"`
		Metadata     []string `json:"metadata"`
	} `json:"scorecard"`
}

// PackageKey is a struct that contains ssystem and name for a package.
type PackageKey struct {
	System string `json:"system"`
	Name   string `json:"name"`
}

// Package Data is a struct that contains PackageKey and Versions for a package.
type PackageData struct {
	PackageKey PackageKey `json:"packageKey"`
	Versions   []Version  `json:"versions"`
}

// Repo is a struct that contains links, versions abd other data about a repository.
type Repo struct {
	Link          string
	AlternateLink string
	ModVersion    string
	OSIData       Project
	VersionData   PackageData
}

// GetOSIData GET call to the OSI API "projects" endpoint.
// See documentation for more detail: https://docs.deps.dev/api/v3alpha/#getproject
func (p *Repo) GetOSIData() error {
	repoName := p.Link

	if len(p.AlternateLink) == 0 && !strings.Contains(p.Link, "github.com") {
		p.AlternateLink = FindGitHub(p.Link)
		repoName = p.AlternateLink
	}

	endpoint := "https://api.deps.dev/v3alpha/projects/" + url.QueryEscape(repoName)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating request for OSI data for %s: %w", repoName, err)
	}

	req = req.WithContext(context.Background())
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error getting OSI data for %s: %w", repoName, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading OSI data for %s: %w", repoName, err)
	}

	err = jsoniter.Unmarshal(body, &p.OSIData)
	if err != nil {
		return fmt.Errorf("error unmarshalling OSI data for %s: %w", repoName, err)
	}

	return nil
}

// GetPackageData GET call to the OSI API "packages" endpoint.
func (p *Repo) GetPackageData() error {
	endpoint := "https://api.deps.dev/v3alpha/systems/go/packages/" + url.QueryEscape(p.Link)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating request for package data for %s: %w", p.Link, err)
	}

	req = req.WithContext(context.Background())
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error getting package data for %s: %w", p.Link, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading package data for %s: %w", p.Link, err)
	}

	err = jsoniter.Unmarshal(body, &p.VersionData)
	if err != nil {
		return fmt.Errorf("error unmarshalling package data for %s: %w", p.Link, err)
	}

	// TODO(unassigned): ensure that the slice of versions are sorted
	return nil
}

// IsLatestVersion returns 1 if the package is the latest version, 0 if it is not, and -1 if there is no OSI data.
func (p *Repo) IsLatestVersion() int {
	if len(p.VersionData.Versions) == 0 {
		return 0
	}

	latestVersion := p.VersionData.Versions[len(p.VersionData.Versions)-1].VersionKey.Version

	return semver.Compare(p.ModVersion, latestVersion)
}

// LatestVersion returns the latest version of the package.
func (p *Repo) LatestVersion() string {
	if len(p.VersionData.Versions) == 0 {
		return "N/A"
	}

	return p.VersionData.Versions[len(p.VersionData.Versions)-1].VersionKey.Version
}

// ExtractRepoInfo takes in a filepath to a go.mod file and returns a slice of Repo structs.
func ExtractRepoInfo(filepath string) []Repo {
	gomod, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return LinksToRepos(gomod)
}

// LinksToRepos takes in a go.mod file as a byte slice and returns a slice of Repo structs.
func LinksToRepos(gomod []byte) []Repo {
	var links []Repo

	lines := bytes.Split(gomod, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("\t")) && !bytes.Contains(line, []byte("indirect")) {
			strLine := string(line)
			links = append(links, Repo{
				Link:       strings.Split(strLine, " ")[0][1:],
				ModVersion: strings.Split(strLine, " ")[1],
				OSIData:    Project{},
			})
		}
	}

	return links
}

// UpdateAllRepos takes in a slice of Repo structs and updates the OSI data for each one.
func UpdateAllRepos(repos []Repo) {
	for i := range repos {
		err := repos[i].GetOSIData()
		if err != nil {
			repos[i].OSIData = Project{}
		}
	}
}
