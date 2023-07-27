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

var mu sync.Mutex


func GoColly(pidsArray []string) {
	log.Printf("Scraping %d items on Amazon, please wait...", len(pidsArray))
	defer log.Printf("Item(s) inserted into database. Waiting for new messages...")
	
	concurrentScrapes := 24
	urlsChannel := make(chan string, len(pidsArray))

	for _, url := range pidsArray {
		urlsChannel <- fmt.Sprintf("https://amazon.com.br/dp/%s", url)
	}
	close(urlsChannel)

	var wg sync.WaitGroup
	wg.Add(concurrentScrapes)

	for i := 0; i < concurrentScrapes; i++ {
		go func(mu *sync.Mutex) {
			defer wg.Done()

			c := colly.NewCollector()
			c.SetRequestTimeout(60 * time.Second)

			c.OnRequest(func(r *colly.Request) {
				fakeHeader := helpers.RandomHeader()
				for key, value := range fakeHeader {
					r.Headers.Set(key, value)
				}
			})

			mu.Lock()
			previousPrice := float32(0)
			elements.GetPreviousPrice(c, &previousPrice)

			price := float32(0)
			elements.GetPrice(c, &price)	
						
			title := "Not Defined"
			elements.GetTitle(c, &title)

			category := ""
			elements.GetCategory(c, &category)

			reviews := 0
			elements.GetReviews(c, &reviews)

			freeShipping := false
			elements.GetFreeShipping(c, &freeShipping)

			imageUrl := "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
			elements.GetImageUrl(c, &imageUrl)
			mu.Unlock()
			
			c.OnScraped(func(r *colly.Response) {
				mu.Lock()
				//previousPrice = elements.ComparePricesAndGetPreviousPrice(price, previousPrice)

				product := types.Product{
					Id: r.Request.URL.String(),
					Title: title,
					Category: category,
					Reviews: reviews,
					Free_Shipping: elements.CheckFreeShipping(freeShipping, category),
					Image_Url: imageUrl,
					Discount: elements.GetDiscount(price, previousPrice),
					Price: price,
					Previous_Price: previousPrice, //elements.GetPreviousPrice(price, discount)
				}
				pg.InsertProduct(pg.UpsertQuery(variables.POSTGRES_PRODUCT_TABLE, product))	
					
				fmt.Printf("\nID: %s\n", product.Id)
				fmt.Printf("TITLE: %s\n", product.Title)
				fmt.Printf("PRICE: %f\n", product.Price)
				fmt.Printf("PREVIOUS: %f\n", product.Previous_Price)
				fmt.Printf("DISCOUNT: %d\n\n", product.Discount)
				mu.Unlock()
			})

			for url := range urlsChannel {
				c.Visit(url)
			}
		}(&mu)
	}
	wg.Wait()
}
