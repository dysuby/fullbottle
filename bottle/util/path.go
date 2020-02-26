package util

import "strings"

func SplitPath(p string) []string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}

	res := strings.Split(p, "/")
	return res
}
