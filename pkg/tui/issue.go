package tui

import (
	"fmt"
	"github.com/knabben/stalker/pkg/testgrid"
	"os"
	"strings"
	"text/template"
)

var basicKeys = []string{"kubetest.Test", "kubetest.DumpClusterLogs", ".Overall"}

func hasBasicTestKeys(name string) bool {
	for _, key := range basicKeys {
		if strings.Contains(name, key) {
			return true
		}
	}
	return false
}

func (d *DashboardIssue) RenderTemplate() {
	fmt.Println(d.renderURL())

	for i, test := range d.Table.Tests {
		if !hasBasicTestKeys(test.Name) {
			if d.Dashboard.OverallStatus == testgrid.FAILING_STATUS {
				d.renderFile()
				fmt.Println(i, test.Name)
			}
		}
	}
	fmt.Println("\n")
}

type Output struct {
	Bla string
}

func (d *DashboardIssue) renderFile() {
	x := Output{Bla: "ble"}
	f, err := template.ParseFiles("pkg/tui/template/failure.tmpl")
	if err != nil {
		panic(err)
	}

	err = f.Execute(os.Stdout, x)
	if err != nil {
		panic(err)
	}
}
