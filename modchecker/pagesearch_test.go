package modchecker

import "testing"

func Test_FindGitHub(t *testing.T) {
	expected := "github.com/go-gorm/mysql"
	actual := FindGitHub("gorm.io/driver/mysql")
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_FindGitHubMultiple(t *testing.T) {
	links := []string{
		"golang.org/x/crypto",
		"gorm.io/driver/mysql",
		"gorm.io/driver/postgres",
		"gorm.io/driver/sqlite",
		"gorm.io/driver/sqlserver",
		"gorm.io/gorm",
	}
	expected := []string{
		"",
		"github.com/go-gorm/mysql",
		"github.com/go-gorm/postgres",
		"github.com/go-gorm/sqlite",
		"github.com/go-gorm/sqlserver",
		"github.com/go-gorm/gorm",
	}
	for i, link := range links {
		actual := FindGitHub(link)
		if expected[i] != actual {
			t.Errorf("expected %v, got %v", expected[i], actual)
		}
	}
}
