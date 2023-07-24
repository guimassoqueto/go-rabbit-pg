package scraper

import (
	"fmt"
	"log"
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
			)

			c := colly.NewCollector()
			c.SetRequestTimeout(60 * time.Second)

			c.OnRequest(func(r *colly.Request) {
				fakeHeader := helpers.RandomHeader()
				for key, value := range fakeHeader {
					r.Headers.Set(key, value)
				}
			})

			c.OnHTML("#title", func(e *colly.HTMLElement) {
				title = strings.ReplaceAll(strings.Trim(e.Text, " "), "'", "''")
			})

			c.OnHTML("#wayfinding-breadcrumbs_container", func(e *colly.HTMLElement) {
				category = elements.GetCategory(e.Text)
			})

			c.OnHTML("#acrCustomerReviewText", func(e *colly.HTMLElement) {
				reviews = elements.GetReviews(e.Text)
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
				}
				pg.InsertProduct(pg.UpsertQuery(variables.POSTGRES_PRODUCT_TABLE, product))
				fmt.Printf("OK: %v", product)
			})

			for url := range urlsChannel {
				c.Visit(url)
			}
		}()
	}
	wg.Wait()
}
