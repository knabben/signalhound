package prow

import "regexp"

func GetRegexParameter(regEx, value string) (params map[string]string) {
	var r = regexp.MustCompile(regEx)
	match := r.FindStringSubmatch(value)
	params = make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			params[name] = match[i]
		}
	}
	return params
}
