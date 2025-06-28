package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/euanwm/modcheck"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func selectModPath() string {
	if len(os.Args) <= 1 {
		if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
			panic("no go.mod file found in current directory")
		}

		return "./go.mod"
	}

	if _, err := os.Stat(os.Args[1] + "go.mod"); os.IsNotExist(err) {
		panic("no go.mod file found in specified directory")
	}

	return os.Args[1] + "go.mod"
}

func setupTable(table *tview.Table) {
	const (
		LinkColumn = iota
		GoModVersionColumn
		LatestVersionColumn
		OpenIssuesColumn
		StarsColumn
		ForksColumn
		OpenSSFScoreColumn
		DescriptionColumn
	)

	table.SetCell(0, LinkColumn, tview.NewTableCell("Link").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, GoModVersionColumn, tview.NewTableCell("GoMod Version").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, LatestVersionColumn, tview.NewTableCell("Latest Version").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, OpenIssuesColumn, tview.NewTableCell("Open Issues").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, StarsColumn, tview.NewTableCell("Stars").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, ForksColumn, tview.NewTableCell("Forks").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, OpenSSFScoreColumn, tview.NewTableCell("OpenSSF Score").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))

	table.SetCell(0, DescriptionColumn, tview.NewTableCell("Description").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignCenter))
}

func populateTable(repos []modcheck.Repo, table *tview.Table) {
	for repo := range repos {
		err := repos[repo].GetPackageData()
		if err != nil {
			panic(err)
		}

		needsUpdating := tcell.ColorBlack

		if repos[repo].IsLatestVersion() > 0 {
			needsUpdating = tcell.ColorRed
		}

		const (
			LinkColumn = iota
			GoModVersionColumn
			LatestVersionColumn
			OpenIssuesColumn
			StarsColumn
			ForksColumn
			OpenSSFScoreColumn
			DescriptionColumn
		)

		table.SetCell(repo+1, LinkColumn, tview.NewTableCell(repos[repo].Link).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, GoModVersionColumn, tview.NewTableCell(repos[repo].ModVersion).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter).SetBackgroundColor(needsUpdating))

		table.SetCell(repo+1, LatestVersionColumn, tview.NewTableCell(repos[repo].LatestVersion()).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, OpenIssuesColumn, tview.NewTableCell(strconv.Itoa(repos[repo].OSIData.OpenIssuesCount)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, StarsColumn, tview.NewTableCell(strconv.Itoa(repos[repo].OSIData.StarsCount)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, ForksColumn, tview.NewTableCell(strconv.Itoa(repos[repo].OSIData.ForksCount)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, OpenSSFScoreColumn, tview.NewTableCell(fmt.Sprint(repos[repo].OSIData.Scorecard.OverallScore)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		table.SetCell(repo+1, DescriptionColumn, tview.NewTableCell(repos[repo].OSIData.Description).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))
	}
}

func main() {
	gomodpath := selectModPath()

	repos := modcheck.ExtractRepoInfo(gomodpath)
	modcheck.UpdateAllRepos(repos)

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	setupTable(table)
	populateTable(repos, table)

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
