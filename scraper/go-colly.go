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


func GoColly(pidsArray []string) {
	log.Printf("Scraping %d items on Amazon, please wait...", len(pidsArray))
	defer log.Printf("Item(s) inserted into database. Waiting for new messages...")
	
	concurrentScrapes := 16
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
				muTitle         sync.Mutex
				muCategory      sync.Mutex
				muReviews       sync.Mutex
				muFreeShipping  sync.Mutex
				muImageUrl      sync.Mutex
				muPrice         sync.Mutex
				muPreviousPrice sync.Mutex
			)

			c := colly.NewCollector()
			c.SetRequestTimeout(60 * time.Second)

			c.OnRequest(func(r *colly.Request) {
				fakeHeader := helpers.RandomHeader()
				for key, value := range fakeHeader {
					r.Headers.Set(key, value)
				}
			})

			muTitle.Lock()
			title := "Not Defined"
			elements.GetTitle(c, &title)
			muTitle.Unlock()

			muCategory.Lock()
			category := ""
			elements.GetCategory(c, &category)
			muCategory.Unlock()

			muReviews.Lock()
			reviews := 0
			elements.GetReviews(c, &reviews)
			muReviews.Unlock()

			muFreeShipping.Lock()
			freeShipping := false
			elements.GetFreeShipping(c, &freeShipping)
			muFreeShipping.Unlock()

			muImageUrl.Lock()
			imageUrl := "https://raw.githubusercontent.com/guimassoqueto/mocks/main/images/404.webp"
			elements.GetImageUrl(c, &imageUrl)
			muImageUrl.Unlock()

			muPrice.Lock()
			price := float32(0)
			elements.GetPrice(c, &price)
			muPrice.Unlock()

			muPreviousPrice.Lock()
			previousPrice := float32(0)
			elements.GetPreviousPrice(c, &previousPrice)
			muPreviousPrice.Unlock()
			

			c.OnScraped(func(r *colly.Response) {
				previousPrice = elements.ComparePricesAndGetPreviousPrice(price, previousPrice)
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
				// fmt.Printf("\nTITLE: %s\n", product.Title)
				// fmt.Printf("ID: %s\n", product.Id)
				// fmt.Printf("FREE SHIPPING: %v\n", product.Free_Shipping)
				// fmt.Printf("IMAGE: %s\n", product.Image_Url)
				// fmt.Printf("CATEGORY: %s\n", product.Category)
				// fmt.Printf("REVIEWS: %d\n", product.Reviews)
				// fmt.Printf("DISCOUNT: %d\n", product.Discount)
				// fmt.Printf("PRICE: %f\n", product.Price)
				// fmt.Printf("PREVIOUS PRICE: %f\n\n", product.Previous_Price)
			})

			for url := range urlsChannel {
				c.Visit(url)
			}
		}()
	}
	wg.Wait()
}
