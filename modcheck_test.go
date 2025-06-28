package modcheck_test

import (
	"testing"

	"github.com/euanwm/modcheck"
)

func Test_extractRepoLinks(t *testing.T) {
	t.Parallel()

	expected := []modcheck.Repo{
		{
			Link:       "github.com/gin-contrib/cors",
			ModVersion: "v1.4.0",
		},
		{
			Link:       "github.com/gin-gonic/gin",
			ModVersion: "v1.9.1",
		},
		{
			Link:       "github.com/heroku/x",
			ModVersion: "v0.0.55",
		},
		{
			Link:       "github.com/json-iterator/go",
			ModVersion: "v1.1.12",
		},
	}
	actual := modcheck.ExtractRepoInfo("./testdata/sample_1/go.mod")

	if len(expected) != len(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_getIOSData(t *testing.T) {
	t.Parallel()

	repo := modcheck.Repo{Link: "github.com/json-iterator/go"}

	err := repo.GetOSIData()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_CheckMod(t *testing.T) {
	t.Parallel()

	modPath := "./testdata/sample_1/go.mod"
	links := modcheck.ExtractRepoInfo(modPath)

	for i := range links {
		err := links[i].GetOSIData()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}
}

func Test_UpdateAllRepos(t *testing.T) {
	t.Parallel()

	modPath := "./testdata/sample_2/go.mod"
	links := modcheck.ExtractRepoInfo(modPath)
	modcheck.UpdateAllRepos(links)

	var count int

	for i := range links {
		if links[i].OSIData.Description != "" {
			count++
		}
	}

	if count != 11 {
		t.Errorf("expected 11, got %v", count)
	}
}

func Test_getPackageData(t *testing.T) {
	t.Parallel()

	repo := modcheck.Repo{Link: "github.com/json-iterator/go"}

	err := repo.GetPackageData()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_CompareVersions(t *testing.T) {
	t.Parallel()

	repo := modcheck.Repo{
		ModVersion: "v3.5.0",
		OSIData:    modcheck.Project{},
		VersionData: modcheck.PackageData{
			Versions: []modcheck.Version{
				{
					VersionKey: modcheck.VersionKey{
						Version: "v3.1.0",
					},
				},
			},
		},
	}

	if repo.IsLatestVersion() != 1 {
		t.Errorf("expected 0 or 1, got %v", repo.IsLatestVersion())
	}
}

func Test_WeirdGormThing(t *testing.T) {
	t.Parallel()

	repo := modcheck.Repo{
		Link:       "gorm.io/gorm",
		ModVersion: "v1.25.1",
	}

	err := repo.GetOSIData()
	if err != nil {
		return
	}

	err = repo.GetPackageData()
	if err != nil {
		return
	}
}
