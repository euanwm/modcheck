package modchecker

import (
	"bytes"
	"github.com/json-iterator/go"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ExtractRepoInfo(filepath string) (links []Repo) {
	gomod, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return linksToRepos(gomod)
}

func linksToRepos(gomod []byte) (links []Repo) {
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

type Repo struct {
	Link          string
	AlternateLink string
	ModVersion    string
	OSIData       Project
	VersionData   PackageData
}

// getOSIData GET call to the OSI API "projects" endpoint.
// See documentation for more detail: https://docs.deps.dev/api/v3alpha/#getproject
func (p *Repo) getOSIData() error {
	repoName := p.Link
	if len(p.AlternateLink) == 0 && !strings.Contains(p.Link, "github.com") {
		p.AlternateLink = FindGitHub(p.Link)
		repoName = p.AlternateLink
	}
	endpoint := "https://api.deps.dev/v3alpha/projects/" + url.QueryEscape(repoName)
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(body, &p.OSIData)
	if err != nil {
		return err
	}
	return nil
}

// Project is the struct that represents the JSON response from the OSI API project endpoint.
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

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackageKey struct {
	System string `json:"system"`
	Name   string `json:"name"`
}

type Version struct {
	VersionKey VersionKey `json:"versionKey"`
	IsDefault  bool       `json:"isDefault"`
}

type PackageData struct {
	PackageKey PackageKey `json:"packageKey"`
	Versions   []Version  `json:"versions"`
}

func (p *Repo) GetPackageData() error {
	endpoint := "https://api.deps.dev/v3alpha/systems/go/packages/" + url.QueryEscape(p.Link)
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(body, &p.VersionData)
	if err != nil {
		return err
	}
	// todo: ensure that the slice of versions are sorted
	return nil
}

func UpdateAllRepos(repos []Repo) {
	for i := range repos {
		err := repos[i].getOSIData()
		if err != nil {
			repos[i].OSIData = Project{}
		}
	}
}

func (p *Repo) IsLatestVersion() int {
	if len(p.VersionData.Versions) == 0 {
		return 0
	}
	latestVersion := p.VersionData.Versions[len(p.VersionData.Versions)-1].VersionKey.Version
	return semver.Compare(p.ModVersion, latestVersion)
}

func (p *Repo) LatestVersion() string {
	if len(p.VersionData.Versions) == 0 {
		return "N/A"
	}
	return p.VersionData.Versions[len(p.VersionData.Versions)-1].VersionKey.Version
}
