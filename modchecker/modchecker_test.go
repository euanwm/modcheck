package modchecker

import (
	"testing"
)

func Test_extractRepoLinks(t *testing.T) {
	expected := []Repo{
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
	actual := ExtractRepoInfo("../test_samples/sample_1/go.mod")
	if len(expected) != len(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_getIOSData(t *testing.T) {
	repo := Repo{Link: "github.com/json-iterator/go"}
	err := repo.getOSIData()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	println(repo.OSIData.Description)
}

func Test_CheckMod(t *testing.T) {
	modPath := "../test_samples/sample_1/go.mod"
	links := ExtractRepoInfo(modPath)
	for i := range links {
		err := links[i].getOSIData()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}
	println(links[0].OSIData.Description)
}

func Test_UpdateAllRepos(t *testing.T) {
	modPath := "../test_samples/sample_2/go.mod"
	links := ExtractRepoInfo(modPath)
	UpdateAllRepos(links)
	var count int
	for i := range links {
		if links[i].OSIData.Description != "" {
			count++
		}
	}
	if count != 10 {
		t.Errorf("expected 10, got %v", count)
	}
}

func Test_getPackageData(t *testing.T) {
	repo := Repo{Link: "github.com/json-iterator/go"}
	err := repo.GetPackageData()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	mostRecentVerions := repo.VersionData.Versions[len(repo.VersionData.Versions)-1].VersionKey.Version
	println(mostRecentVerions)
}

func Test_CompareVersions(t *testing.T) {
	repo := Repo{
		ModVersion: "v3.5.0",
		OSIData:    Project{},
		VersionData: PackageData{
			Versions: []Version{
				{
					VersionKey: VersionKey{
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
	repo := Repo{
		Link:       "gorm.io/gorm",
		ModVersion: "v1.25.1",
	}
	err := repo.getOSIData()
	if err != nil {
		return
	}
	err = repo.GetPackageData()
	if err != nil {
		return
	}
	println(repo.VersionData.Versions[len(repo.VersionData.Versions)-1].VersionKey.Version)
}
