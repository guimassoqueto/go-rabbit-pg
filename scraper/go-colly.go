package scraper

import (
	"fmt"
	"log"
	"sync"
	"time"

	"grp/elements"
	"grp/helpers"
	pg "grp/postgres"
	"grp/types"
	"grp/variables"

	"github.com/gocolly/colly"
)

func goColly(urlCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urlCh {
		var (
			title         string
			category      string
			reviews       int
			freeShipping  bool
			imageUrl      string = "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
			price         float32
			previousPrice float32
		)

		c := colly.NewCollector()
		c.SetRequestTimeout(60 * time.Second)

		c.OnRequest(func(r *colly.Request) {
			fakeHeader := helpers.RandomHeader()
			for key, value := range fakeHeader {
				r.Headers.Set(key, value)
			}
		})

		elements.GetPreviousPrice(c, &previousPrice)
		elements.GetPrice(c, &price)
		elements.GetFreeShipping(c, &freeShipping)
		elements.GetReviews(c, &reviews)
		elements.GetImageUrl(c, &imageUrl)
		elements.GetTitle(c, &title)
		elements.GetCategory(c, &category)

		c.OnScraped(func(r *colly.Response) {
			previousPrice = elements.ComparePricesAndGetPreviousPrice(price, previousPrice)
			product := types.Product{
				Id:             r.Request.URL.String(),
				Title:          title,
				Category:       category,
				Reviews:        reviews,
				Free_Shipping:  elements.CheckFreeShipping(freeShipping, category),
				Image_Url:      imageUrl,
				Discount:       elements.GetDiscount(price, previousPrice),
				Price:          price,
				Previous_Price: previousPrice,
			}
			if product.Price > 0 {
				pg.InsertProduct(pg.UpsertQuery(variables.POSTGRES_TABLE_NAME, product))
			}
		})
		c.Visit(url)
	}
}

func Scrap(pidsArray []string) {
	log.Printf("Scraping %d items on Amazon, please wait...", len(pidsArray))
	defer log.Printf("Items insert into database. Waiting for new pids...")

	var maxConcurrentRequests = variables.MAX_CONCURRENCY()
	var wg sync.WaitGroup

	urlCh := make(chan string)

	for i := 0; i < maxConcurrentRequests; i++ {
		wg.Add(1)
		go goColly(urlCh, &wg)
	}

	for _, pid := range pidsArray {
		urlCh <- fmt.Sprintf("https://amazon.com.br/dp/%s", pid)
	}

	close(urlCh)
	wg.Wait()
}
