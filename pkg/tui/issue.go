package tui

import (
	"fmt"
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

func (d *DashboardTab) RenderTemplate() {
	fmt.Println(d.renderURL())
}

type Output struct {
	Bla string
}

func (d *DashboardTab) renderFile() {
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
