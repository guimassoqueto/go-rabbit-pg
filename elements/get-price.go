package elements

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)


func priceToFloat(eText string) float32 {
	regex, _ := regexp.Compile(`[\d\,\.]+`)
	match := regex.FindString(eText)
	removedDots := strings.ReplaceAll(match, ".", "")
	commaToDot := strings.ReplaceAll(removedDots, ",", ".")
	floatPrice, error := strconv.ParseFloat(commaToDot, 32)
	if error != nil {
		return 0
	}
	return float32(floatPrice)
}


func GetPrice(c *colly.Collector, price *float32) {
	// PADRÃO E TABELA
	var temp float32
	c.OnHTML("#corePrice_feature_div", func(e *colly.HTMLElement) {
		temp = priceToFloat(e.Text)
		*price = temp
	})
	// EM CERTOS CASOS A #corePrice_feature_div POSSUI ELEMENTOS ADICIONAIS
	c.OnHTML("#corePrice_feature_div>div>div>span>span.a-offscreen", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*price = temp
		}
	})

	// KINDLE
	c.OnHTML("#kindle-price", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*price = temp
		}
	})
	// LIVRO FÍSICO
	c.OnHTML("#price", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*price = temp
		}
	})
	// OCULOS SPEEDO = B07FW11H5X
	c.OnHTML(".apexPriceToPay>span", func(e *colly.HTMLElement) {
		if temp == 0 {
			temp = priceToFloat(e.Text)
			*price = temp
		}
	})
}