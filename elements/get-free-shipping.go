package elements

import (
	"strings"

	"github.com/gocolly/colly"
)

func CheckFreeShipping(freeShipping bool, category string) bool {
	if strings.Contains(strings.ToLower(category), "livro") {
		return true
	}
	if strings.Contains(strings.ToLower(category), "kindle") {
		return true
	}
	if strings.Contains(strings.ToLower(category), "ebook") {
		return true
	}
	return freeShipping
}

func GetFreeShipping(c *colly.Collector, freeShipping *bool) {
	c.OnHTML("#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE>span", func(e *colly.HTMLElement) {
		attr := strings.Trim(e.Attr("data-csa-c-delivery-price"), " ")
		if attr == "GRÁTIS" { 
			*freeShipping = true 
		} else {
			*freeShipping = false
		}
	})
}