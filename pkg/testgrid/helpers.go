package testgrid

import "strings"

func CleanSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "%20")
}
