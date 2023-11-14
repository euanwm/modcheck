package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"modchecker/modchecker"
	"os"
)

func main() {
	var gomodpath string
	if len(os.Args) <= 1 {
		if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
			panic("no go.mod file found in current directory")
		}
		gomodpath = "./go.mod"
	} else {
		if _, err := os.Stat(os.Args[1] + "go.mod"); os.IsNotExist(err) {
			panic("no go.mod file found in specified directory")
		}
		gomodpath = os.Args[1] + "go.mod"
	}

	repos := modchecker.ExtractRepoInfo(gomodpath)
	modchecker.UpdateAllRepos(repos)

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	table.SetCell(0, 0, tview.NewTableCell("Link").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("GoMod Version").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Latest Version").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Open Issues").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Stars").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Forks").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 6, tview.NewTableCell("OpenSSF Score").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	table.SetCell(0, 7, tview.NewTableCell("Description").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))

	for i := range repos {
		err := repos[i].GetPackageData()
		if err != nil {
			return
		}
		needsUpdating := tcell.ColorBlack

		if repos[i].IsLatestVersion() > 0 {
			needsUpdating = tcell.ColorRed
		}

		table.SetCell(i+1, 0, tview.NewTableCell(repos[i].Link).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(repos[i].ModVersion).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter).SetBackgroundColor(needsUpdating))
		table.SetCell(i+1, 2, tview.NewTableCell(repos[i].LatestVersion()).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 3, tview.NewTableCell(repos[i].OSIData.OpenIssuesCount).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 4, tview.NewTableCell(repos[i].OSIData.StarsCount).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 5, tview.NewTableCell(repos[i].OSIData.ForksCount).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 6, tview.NewTableCell(fmt.Sprint(repos[i].OSIData.Scorecard.OverallScore)).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 7, tview.NewTableCell(repos[i].OSIData.Description).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	}

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
