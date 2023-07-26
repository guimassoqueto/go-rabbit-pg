package elements

import (
	"grp/helpers"

	"github.com/gocolly/colly"
)

func getBiggestImage(e *colly.HTMLElement, attr string) string {
	var imageUrl string = "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
	if e.Attr(attr) != "" {
		var availableImages []string = helpers.StringfiedJsonKeysToArray(e.Attr(attr))
		imageUrl = availableImages[len(availableImages)-1]
	}
	return imageUrl
}

func GetImageUrl(c *colly.Collector, imageUrl *string) {
	c.OnHTML("#ebooksImgBlkFront", func(e *colly.HTMLElement) {
		*imageUrl = getBiggestImage(e, "data-a-dynamic-image")
	})
	c.OnHTML("#imgBlkFront", func(e *colly.HTMLElement) {
		*imageUrl = getBiggestImage(e, "data-a-dynamic-image")
	})
	c.OnHTML("img.a-dynamic-image", func(e *colly.HTMLElement) {
		if e.Attr("data-old-hires") != "" {
			*imageUrl = e.Attr("data-old-hires")
		}
	})
	c.OnHTML("#landingImage", func(e *colly.HTMLElement) {
		if e.Attr("data-old-hires") != "" {
			*imageUrl = e.Attr("data-old-hires")
		} else {
			*imageUrl = getBiggestImage(e, "data-a-dynamic-image")
		}
	})
}
