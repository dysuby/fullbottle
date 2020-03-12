package util

import "strings"

func SplitPath(p string) []string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}

	res := strings.Split(p, "/")
	if res[0] == "" {
		res = res[1:]
	}
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}
	return res
}
