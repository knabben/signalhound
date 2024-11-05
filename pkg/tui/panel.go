package tui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/knabben/stalker/pkg/github"
	"github.com/knabben/stalker/pkg/testgrid"
	"github.com/rivo/tview"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os/exec"
	"strings"
	"time"
)

var (
	pagesName   = "Stalker"
	app         *tview.Application // The tview application.
	pages       *tview.Pages       // The application pages.
	brokenPanel = tview.NewList()
	position    = tview.NewTextView()
	slackPanel  = tview.NewTextArea()
	githubPanel = tview.NewTextArea()
)

var positionText = "[yellow]Select a content Windows and press [blue]Ctrl-Space [yellow]to COPY or press [blue]Ctrl-C [yellow]to exit"

func RenderVisual(tabs []*DashboardTab, githubToken string) error {
	app = tview.NewApplication()

	// Render tab
	tabsPanel := tview.NewList().ShowSecondaryText(false)
	tabsPanel.SetBorder(true).SetTitle("Board#Tabs")

	// Final issue rendering
	slackPanel.SetBorder(true).SetTitle("Slack Message")
	slackPanel.SetWrap(true).SetDisabled(true)
	githubPanel.SetBorder(true).SetTitle("Github Issue")
	githubPanel.SetWrap(true)

	// Broken tests in the tab
	brokenPanel.ShowSecondaryText(false).SetDoneFunc(func() { app.SetFocus(tabsPanel) })
	brokenPanel.SetBorder(true).SetTitle("Tests")

	position.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText(positionText)

	// Create the layout.
	grid := tview.NewGrid().SetRows(13, 13, 0, 1).
		AddItem(tabsPanel, 0, 0, 1, 2, 0, 0, true).
		AddItem(brokenPanel, 1, 0, 1, 2, 0, 0, false).
		AddItem(slackPanel, 2, 0, 1, 1, 0, 0, false).
		AddItem(githubPanel, 2, 1, 1, 1, 0, 0, false).
		AddItem(position, 3, 0, 1, 2, 0, 0, false)

	for _, tab := range tabs {
		tabsPanel.AddItem(tab.BoardHash, "", 0, func() {
			brokenPanel.Clear()
			for _, test := range tab.Tests {
				brokenPanel.AddItem(test.Name, "", 0, nil)
			}
			app.SetFocus(brokenPanel)
			brokenPanel.SetCurrentItem(0)
			brokenPanel.SetChangedFunc(func(i int, testName string, t string, s rune) {
				position.SetText(positionText)
				slackPanel.SetBorderColor(tcell.ColorWhite)
				githubPanel.SetBorderColor(tcell.ColorWhite)
			})
			brokenPanel.SetSelectedFunc(func(i int, testName string, t string, s rune) {
				var currentTest = tab.Tests[i]
				updateSlackPanel(tab, currentTest)
				updateGitHubPanel(tab, currentTest, githubToken)
				app.SetFocus(slackPanel)
			})
		})
	}

	pages = tview.NewPages().AddPage(pagesName, grid, true, true)
	return app.SetRoot(pages, true).EnableMouse(true).Run()
}

// updateSlackPanel writes down to left panel (Slack) content.
func updateSlackPanel(tab *DashboardTab, currentTest *TabTest) {
	// set the item string with current test content
	item := fmt.Sprintf("%s %s on [%s](%s): `%s` [Prow](%s), [Triage](%s), last failure on %s\n",
		tab.Icon, cases.Title(language.English).String(tab.State), tab.BoardHash, tab.BoardURL,
		currentTest.Name, currentTest.ProwURL, currentTest.TriageURL, timeClean(currentTest.LatestTimestamp),
	)

	// set input capture, ctrl-space for clipboard copy, esc to cancel panel selection.
	slackPanel.SetText(item, true)
	slackPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			position.SetText("[blue]COPIED [yellow]SLACK [blue]TO THE CLIPBOARD!")
			if err := CopyToClipboard(slackPanel.GetText()); err != nil {
				position.SetText(fmt.Sprintf("[red]error: %v", err.Error()))
				return event
			}
			slackPanel.SetBorderColor(tcell.ColorBlue)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyUp {
			slackPanel.SetText("", false)
			githubPanel.SetText("", false)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyRight {
			app.SetFocus(githubPanel)
		}
		return event
	})
}

// updateGitHubPanel writes down to the right panel (GitHub) content.
func updateGitHubPanel(tab *DashboardTab, currentTest *TabTest, token string) {
	// create the filled out issue template object
	splitBoard := strings.Split(tab.BoardHash, "#")
	issue := &IssueTemplate{
		BoardName:    splitBoard[0],
		TabName:      splitBoard[1],
		TestName:     currentTest.Name,
		TestGridURL:  tab.BoardURL,
		TriageURL:    currentTest.TriageURL,
		ProwURL:      currentTest.ProwURL,
		ErrMessage:   currentTest.ErrMessage,
		FirstFailure: timeClean(currentTest.FirstTimestamp),
		LastFailure:  timeClean(currentTest.LatestTimestamp),
	}

	// pick the correct template by failure status
	templateFile, prefixTitle := "template/flake.tmpl", "Flaking Test"
	if tab.Status == testgrid.FAILING_STATUS {
		templateFile, prefixTitle = "template/failure.tmpl", "Failing Test"
	}
	template, err := tab.renderTemplate(issue, templateFile)
	if err != nil {
		position.SetText(fmt.Sprintf("[red]error: %v", err.Error()))
		return
	}
	issueTemplate := template.String()
	issueTitle := fmt.Sprintf("[%v] %v", prefixTitle, currentTest.Name)
	githubPanel.SetText(issueTemplate, false)

	// set input capture, ctrl-space for clipboard copy, ctrl-b for
	// automatic GitHub draft issue creation.
	githubPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			position.SetText("[blue]COPIED [yellow]ISSUE [blue]TO THE CLIPBOARD!")
			if err := CopyToClipboard(githubPanel.GetText()); err != nil {
				position.SetText(fmt.Sprintf("[red]error: %v", err.Error()))
				return event
			}
			githubPanel.SetBorderColor(tcell.ColorBlue)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyCtrlB {
			gh := github.NewGithub(context.Background(), token)
			if err := gh.CreateDraftIssue(issueTitle, issueTemplate); err != nil {
				position.SetText(fmt.Sprintf("[red]error: %v", err.Error()))
				return event
			}
			position.SetText("[blue]Created [yellow]DRAFT ISSUE [blue] on GitHub Project!")
			githubPanel.SetBorderColor(tcell.ColorBlue)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyEscape {
			slackPanel.SetText("", false)
			githubPanel.SetText("", false)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(slackPanel)
		}
		return event
	})
}

// timeClean returns the string representation of the timestamp.
func timeClean(ts int64) string {
	return time.Unix(ts/1000, 0).Format(time.RFC1123)
}

// CopyToClipboard pipes the panel content to clip.exe WSL.
func CopyToClipboard(text string) error {
	args := "echo '" + text + "' | clip.exe"
	cmd := exec.Command("bash", "-c", args)
	return cmd.Run()
}
