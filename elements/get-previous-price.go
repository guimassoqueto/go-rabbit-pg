package elements

import (
	"github.com/gocolly/colly"
)


func ComparePricesAndGetPreviousPrice(price float32, previousPrice float32) float32 {
	if price > previousPrice  {
		return price
	}
	return previousPrice
}

func GetPreviousPrice(c *colly.Collector, previousPrice *float32) {
	firstOccurrenceProcessed := false
	// PADRAO
	c.OnHTML(".basisPrice", func(e *colly.HTMLElement) {
		subElement := e.DOM.Find("span.a-offscreen").First()
		if !firstOccurrenceProcessed {
			*previousPrice = priceToFloat(subElement.Text())
			firstOccurrenceProcessed = true
		}
	})

	// EBOOKS
	c.OnHTML("#digital-list-price", func(e *colly.HTMLElement) {
		*previousPrice = priceToFloat(e.Text)
	})

	// TABELA
	c.OnHTML("#corePrice_desktop", func(e *colly.HTMLElement) {
		subElement := e.DOM.Find("span.a-offscreen").First()
		if !firstOccurrenceProcessed {
			*previousPrice = priceToFloat(subElement.Text())
			firstOccurrenceProcessed = true
		}
	})

	// LIVROFíSICO
	c.OnHTML("#listPrice", func(e *colly.HTMLElement) {
		*previousPrice = priceToFloat(e.Text)
	})
}