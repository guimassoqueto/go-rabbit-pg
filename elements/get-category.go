package elements

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)


func extractCategory(eText string) string {
	regex, _ := regexp.Compile(`[a-zA-Zà-úÀ-Ú0-9_]+`)
	match := regex.FindAllString(eText, -1)
	return strings.Join(match, " ")
}

func GetCategory(c *colly.Collector, category *string) {
	c.OnHTML("#wayfinding-breadcrumbs_container", func(e *colly.HTMLElement) {
		*category = extractCategory(e.Text)
	})
}
