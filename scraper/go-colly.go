package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"grp/elements"
	"grp/helpers"
	pg "grp/postgres"
	"grp/types"
	"grp/variables"

	"github.com/gocolly/colly"
)

func getBiggestImage(e *colly.HTMLElement, attr string) string {
	var imageUrl string = "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
	if e.Attr(attr) != "" {
		var availableImages []string = helpers.StringfiedJsonKeysToArray(e.Attr(attr))
		imageUrl = availableImages[len(availableImages) - 1]
	}
	return imageUrl
}

func convertDiscountToInteger(eText string) int {
	regex, _ := regexp.Compile(`\d+%`)
	match := regex.FindString(eText)
	percentageRemoved := strings.ReplaceAll(match, "%", "")
	discount, err := strconv.Atoi(percentageRemoved)
	if err != nil {
		return 0
	}
	return discount
}

func GoColly(pidsArray []string) {
	log.Printf("Scraping %d items on Amazon, please wait...", len(pidsArray))
	defer log.Printf("Item(s) inserted into database. Waiting for new messages...")
	
	concurrentScrapes := 32
	urlsChannel := make(chan string, len(pidsArray))

	for _, url := range pidsArray {
		urlsChannel <- fmt.Sprintf("https://amazon.com.br/dp/%s", url)
	}
	close(urlsChannel)

	var wg sync.WaitGroup
	wg.Add(concurrentScrapes)

	for i := 0; i < concurrentScrapes; i++ {
		go func() {
			defer wg.Done()

			var (
				title string = "Not Defined"
				category string = "Not Definded"
				reviews int = 0
				freeShipping bool = false
				imageUrl string = "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
				discount int = 0
			)

			c := colly.NewCollector()
			c.SetRequestTimeout(60 * time.Second)

			c.OnRequest(func(r *colly.Request) {
				fakeHeader := helpers.RandomHeader()
				for key, value := range fakeHeader {
					r.Headers.Set(key, value)
				}
			})
			// TITLE
			c.OnHTML("#title", func(e *colly.HTMLElement) {
				title = strings.ReplaceAll(strings.Trim(e.Text, " "), "'", "''")
			})
			// CATEGORY
			c.OnHTML("#wayfinding-breadcrumbs_container", func(e *colly.HTMLElement) {
				category = elements.GetCategory(e.Text)
			})
			// REVIEWS
			c.OnHTML("#acrCustomerReviewText", func(e *colly.HTMLElement) {
				reviews = elements.GetReviews(e.Text)
			})
			// FREE-SHIPPING
			c.OnHTML("#primeSavingsUpsellCaption_feature_div", func(e *colly.HTMLElement) { freeShipping = true })
			c.OnHTML("div.tabular-buybox-text:nth-child(4)>div:nth-child(1)>span:nth-child(1)", func(e *colly.HTMLElement) {
				innerText := strings.ToLower(e.Text)
				if strings.Contains(innerText, "amazon") {
					freeShipping = true
				}
			})
			c.OnHTML("div#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE>span", func(e *colly.HTMLElement) {
				innerText := strings.ToLower(e.Text)
				if strings.Contains(innerText, "grátis") {
					freeShipping = true
				}
			})
			// IMAGE-URL
			c.OnHTML("img.a-dynamic-image", func(e *colly.HTMLElement) {
				if e.Attr("data-old-hires") != "" {
					imageUrl = e.Attr("data-old-hires")
				}
			})
			c.OnHTML("#landingImage", func(e *colly.HTMLElement) {
				if e.Attr("data-old-hires") != "" {
					imageUrl = e.Attr("data-old-hires")
				} else {
					imageUrl = getBiggestImage(e, "data-a-dynamic-image")
				}
			})
			c.OnHTML("#ebooksImgBlkFront", func(e *colly.HTMLElement) {
				imageUrl = getBiggestImage(e, "data-a-dynamic-image")
			})
			c.OnHTML("#imgBlkFront", func(e *colly.HTMLElement) {
				imageUrl = getBiggestImage(e, "data-a-dynamic-image")
			})
			// DISCOUNT
			c.OnHTML(".savingPriceOverride", func(e *colly.HTMLElement) {
				discount = convertDiscountToInteger(e.Text)
			})
			c.OnHTML("#savingsPercentage", func(e *colly.HTMLElement) {
				discount = convertDiscountToInteger(e.Text)
			})
			c.OnHTML("p.ebooks-price-savings", func(e *colly.HTMLElement) {
				discount = convertDiscountToInteger(e.Text)
			})
			c.OnHTML("tr>td.a-span12.a-color-price.a-size-base>span.a-color-price", func(e *colly.HTMLElement) {
				discount = convertDiscountToInteger(e.Text)
			})

			c.OnError(func(r *colly.Response, err error) {
				log.Printf("Error while scraping an item: %s", err.Error())
			})
			
			c.OnScraped(func(r *colly.Response) {
				product := types.Product{
					Id: r.Request.URL.String(),
					Title: title,
					Category: category,
					Reviews: reviews,
					Free_Shipping: freeShipping,
					Image_Url: imageUrl,
					Discount: discount,
				}
				pg.InsertProduct(pg.UpsertQuery(variables.POSTGRES_PRODUCT_TABLE, product))
				log.Printf("OK: %v", product)
			})

			for url := range urlsChannel {
				c.Visit(url)
			}
		}()
	}
	wg.Wait()
}
