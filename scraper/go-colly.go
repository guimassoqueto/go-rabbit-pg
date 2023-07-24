package scraper

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	//"grp/postgres"
	"grp/types"
	//"grp/variables"

	"github.com/gocolly/colly"
	randomHeader "github.com/guimassoqueto/go-fake-headers"
)

func GoColly(pidsArray []string) {
	var (
		title string
	)

	var wg sync.WaitGroup

	concurrentLinks := make(chan string, 32)

	c := colly.NewCollector(
		colly.AllowedDomains(),
		colly.IgnoreRobotsTxt(),
	)

	c.SetRequestTimeout(60 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		fakeHeader := randomHeader.Build()
		for key, value := range fakeHeader {
			r.Headers.Set(key, value)
		}
	})

	c.OnHTML("#title", func(e *colly.HTMLElement) {
		title = strings.ReplaceAll(strings.Trim(e.Text, " "), "'", "''")
		wg.Done()
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		wg.Done()
	})

	c.OnScraped(func(r *colly.Response) {
		product := types.Product{
			Id:    r.Request.URL.String(),
			Title: title,
		}
		// pg.InsertProduct(pg.UpsertQuery(variables.POSTGRES_PRODUCT_TABLE, product))
		log.Printf("OK: %s", product)
	})

	// Start Goroutines for each URL
	for _, url := range pidsArray {
		wg.Add(1) // Increment the WaitGroup counter for each URL

		// Send the URL to the worker pool for processing
		concurrentLinks <- url

		go func(u string) {
			// Defer the removal of the URL from the worker pool
			defer func() { <-concurrentLinks }()

			// Make the request using Colly
			err := c.Visit(u)
			if err != nil {
				log.Println("Error visiting URL:", u, "\nError:", err)
				return
			}
		}(fmt.Sprintf("https://amazon.com.br/dp/%s", url))
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(concurrentLinks)
}
