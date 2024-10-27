package tui

import (
	"fmt"
	"strings"
)

func (d *DashboardIssue) RenderTemplate() {
	fmt.Println(d.renderURL())

	for i, test := range d.Table.Tests {
		if !hasBasicTestKeys(test.Name) {
			fmt.Println(i, test.Name)
		}
	}
	fmt.Println("\n")
}

var basicKeys = []string{"kubetest.Test", "kubetest.DumpClusterLogs", ".Overall"}

func hasBasicTestKeys(name string) bool {
	for _, key := range basicKeys {
		if strings.Contains(name, key) {
			return true
		}
	}
	return false

}
