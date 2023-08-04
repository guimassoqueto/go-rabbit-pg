package elements

import (
	"strings"
	"github.com/gocolly/colly"
)

func GetTitle(c *colly.Collector, title *string) {
	c.OnHTML("#productTitle", func(e *colly.HTMLElement) {
		removedSingleQuotes := strings.ReplaceAll(strings.Trim(e.Text, " "), "'", "''")
		*title = strings.ReplaceAll(removedSingleQuotes, "\"", "")
	})
}