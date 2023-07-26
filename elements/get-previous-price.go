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
	var temp float32 = 0

	// PADRÃO
	c.OnHTML(".basisPrice", func(e *colly.HTMLElement) {
		temp = priceToFloat(e.Text)
		*previousPrice = temp
	})
	// EBOOKS
	c.OnHTML("#digital-list-price", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*previousPrice = temp
		}
	})
	// TABELA
	c.OnHTML("tbody>tr>td>span>span.a-offscreen", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*previousPrice = temp
		}
	})
	// LIVROFíSICO
	c.OnHTML("#listPrice", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*previousPrice = temp
		}
	})
}