package elements

import (
	"regexp"
	"strings"
)

func GetCategory(eText string) string {
	regex, _ := regexp.Compile(`[a-zA-Zà-úÀ-Ú0-9_]+`)
	match := regex.FindAllString(eText, -1)
	return strings.Join(match, " ")
}
