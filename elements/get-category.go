package elements

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func GetCategory(e *colly.HTMLElement) string {
	regex, _ := regexp.Compile(`[a-zA-Zà-úÀ-Ú0-9_]+`)
	match := regex.FindAllString(e.Text, -1)
	return strings.Join(match, " ")
}
