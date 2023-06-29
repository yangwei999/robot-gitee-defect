package utils

import "strings"

func TrimString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return strings.Replace(str, "\t", "", -1)
}
