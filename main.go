package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	Image  string `json:"image"`
	Rating string `json:"rating"`
}

func main() {
        var keyword string

        for {
                fmt.Print("Please input the keyword you wanna scrap for: ")
        fmt.Scanln(&keyword)
                if (len([]rune(keyword)) != 0) {
                    break
                }
            }

	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	var products []Product

	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(h *colly.HTMLElement) {
		h.ForEach("div.a-section.a-spacing-base", func(_ int, h *colly.HTMLElement) {
			product := Product{
				Name:   h.ChildText("span.a-size-base-plus"),
				Image:  h.ChildAttr("img.s-image", "src"),
				Price:  h.ChildText("span.a-price-whole"),
				Rating: h.ChildText("span.a-icon-alt"),
			}

			products = append(products, product)
			fmt.Println(product)
		})
	})

	c.OnHTML("a.s-pagination-item.s-pagination-next.s-pagination-button.s-pagination-separator", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Link of the page:", r.URL.String())
	})

	c.Visit("https://www.amazon.in/s?k=" + keyword)
	fmt.Println(products)
	content, err := json.Marshal(products)

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("products.json", content, 0644)
}
