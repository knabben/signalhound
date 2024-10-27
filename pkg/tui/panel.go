package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os/exec"
	"time"
)

var (
	pagesName = "Stalker"
	app       *tview.Application // The tview application.
	pages     *tview.Pages       // The application pages.
	focus     tview.Primitive    // The primitive in the focus that last had focus.

	brokenPanel = tview.NewList()
	position    = tview.NewTextView()
	renderPanel = tview.NewTextArea()
)

var positionText = "[yellow]Select an option and press Ctrl-Space to COPY or [white]press Ctrl-C to exit"

func RenderVisual(tabs []*DashboardTab) error {
	app = tview.NewApplication()

	// Render tab
	tabsPanel := tview.NewList().ShowSecondaryText(false)
	tabsPanel.SetBorder(true).SetTitle("Tabs")

	// Final issue rendering
	renderPanel.SetWrap(true).SetBorder(true).SetTitle("Message")
	renderPanel.SetDisabled(true)

	// Broken tests in the tab
	brokenPanel.ShowSecondaryText(false).SetDoneFunc(func() { app.SetFocus(tabsPanel) })
	brokenPanel.SetBorder(true).SetTitle("Tests")

	position.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText(positionText)

	// Create the layout.
	grid := tview.NewGrid().SetRows(15, 15, 0, 1).
		AddItem(tabsPanel, 0, 0, 1, 1, 0, 0, true).
		AddItem(brokenPanel, 1, 0, 1, 1, 0, 0, false).
		AddItem(renderPanel, 2, 0, 1, 1, 0, 0, false).
		AddItem(position, 3, 0, 1, 1, 0, 0, false)

	for _, tab := range tabs {
		tabsPanel.AddItem(tab.BoardHash, "", 0, func() {
			brokenPanel.Clear()
			for _, test := range tab.Tests {
				brokenPanel.AddItem(test.Name, "", 0, nil)
			}
			app.SetFocus(brokenPanel)
			brokenPanel.SetCurrentItem(0)
			brokenPanel.SetSelectedFunc(func(i int, testName string, t string, s rune) {
				contentRender(tab, i)
				app.SetFocus(renderPanel)
			})
		})
	}

	pages = tview.NewPages().AddPage(pagesName, grid, true, true)
	return app.SetRoot(pages, true).EnableMouse(true).Run()
}

func contentRender(tab *DashboardTab, i int) {
	currentTest := tab.Tests[i]
	unixTimeUTC := time.Unix(currentTest.LatestTimestamp/1000, 0)

	item := fmt.Sprintf("%s %s on [%s](%s): `%s` [Prow](%s), [Triage](%s), last failure on %s\n",
		tab.Icon, tab.State, tab.BoardHash, tab.BoardURL, currentTest.Name, currentTest.ProwURL, currentTest.TriageURL, unixTimeUTC.Format(time.RFC822))

	// write down to the panel
	renderPanel.SetText(item, true)
	renderPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			position.SetText("[blue]COPIED TO THE CLIPBOARD!")
			args := "echo '" + renderPanel.GetText() + "' | clip.exe"
			cmd := exec.Command("bash", "-c", args)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			renderPanel.SetText("", true)
			app.SetFocus(brokenPanel)
		}
		if event.Key() == tcell.KeyEscape {
			renderPanel.SetText("", true)
			app.SetFocus(brokenPanel)
		}
		return event
	})
}
