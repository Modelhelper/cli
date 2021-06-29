package slice

import "strings"

func MaxLen(list []string) int {

	l := 0

	for _, v := range list {
		tl := len(v)
		if l < tl {
			l = tl
		}
	}

	return l
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}

	return false
}
