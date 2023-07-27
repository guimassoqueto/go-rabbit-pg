package elements

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

var mu sync.Mutex


func priceToFloat(eText string) float32 {
	regex, _ := regexp.Compile(`[\d\,\.]+`)
	match := regex.FindString(eText)
	if strings.Contains(match, ".") {
		match = strings.ReplaceAll(match, ".", "")
	}
	commaToDot := strings.ReplaceAll(match, ",", ".")
	floatPrice, error := strconv.ParseFloat(commaToDot, 32)

	if error != nil {
		return 0
	}
	return float32(floatPrice)
}


func GetPrice(c *colly.Collector, price *float32) {
	// PADRÃO E TABELA
	firstOccurrenceProcessed := false
	c.OnHTML("#corePrice_feature_div", func(e *colly.HTMLElement) {
		subElement := e.DOM.Find("span.a-offscreen").First()
		if !firstOccurrenceProcessed {
			*price = priceToFloat(subElement.Text())
			firstOccurrenceProcessed = true
		}
	})

	// KINDLE
	c.OnHTML("#kindle-price", func(e *colly.HTMLElement) {
		*price = priceToFloat(e.Text)
	})
	// LIVRO FÍSICO
	c.OnHTML("#price", func(e *colly.HTMLElement) {
		*price = priceToFloat(e.Text)
	})
}