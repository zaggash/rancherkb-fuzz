package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	md "github.com/JohannesKaufmann/html-to-markdown"
	colly "github.com/gocolly/colly/v2"
)

type Article struct {
	Id      string
	Title   string
	Url     string
	Content string
}

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.URLFilters(
			regexp.MustCompile("www\\.suse\\.com/support/kb/(.*)/?id(.*)"),
		),
	)

	articlesCollector := c.Clone()
	articles := make([]Article, 0)

	c.OnHTML(".result-table", func(e *colly.HTMLElement) {
		e.ForEach(".result-cell a[href]", func(i int, h *colly.HTMLElement) {
			articlesCollector.Visit(h.Request.AbsoluteURL(h.Attr("href")))
		})
	})

	articlesCollector.OnHTML(".col_one", func(e *colly.HTMLElement) {
		item := Article{}
		url := e.Request.URL.String()
		title := e.ChildText("h1")
		content, _ := e.DOM.Html()

		item.Id = url[len(url)-9:]
		item.Url = url
		item.Title = title
		item.Content = content
		articles = append(articles, item)

		fmt.Println("-------")
		fmt.Println("Id:", item.Id)
		fmt.Println("Link:", item.Url)
		fmt.Println("Title:", item.Title)
		//fmt.Println("Content:", item.Content)

		converter := md.NewConverter("", true, nil)
		markdown, err := converter.ConvertString(item.Content)

		file, err := os.Create("./website/docs/" + item.Id + ".md")
		if err != nil {
			fmt.Println(err)
		} else {
			file.WriteString(markdown)
		}
		file.Close()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	articlesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Article", r.URL.String())
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnHTML(".results_summary a[href]", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})
	articlesCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(articles, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("articles.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	c.Visit("https://www.suse.com/support/kb/?id=SUSE+Rancher")
}
