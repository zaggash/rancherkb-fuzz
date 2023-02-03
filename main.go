package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	colly "github.com/gocolly/colly/v2"
)

type Article struct {
	Id      string
	Title   string
	Url     string
	Content string
}

func main() {
	book := make([]Article, 0)

	titleCollector := colly.NewCollector(
		colly.MaxDepth(1),
		colly.URLFilters(
			regexp.MustCompile(`www.suse.com/support/kb/(.*)/?id(.*)`),
		),
	)

	articleCollector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
	)
	articleCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})

	titleCollector.OnHTML(".result-table", func(titleHtmlTable *colly.HTMLElement) {
		titleHtmlTable.ForEach(".result-cell a[href]", func(i int, titleHtmlLink *colly.HTMLElement) {
			articleCollector.Visit(titleHtmlLink.Request.AbsoluteURL(titleHtmlLink.Attr("href")))
		})
		articleCollector.Wait()
	})

	articleCollector.OnHTML(".col_one", func(mainHtmlArticle *colly.HTMLElement) {
		page := Article{}
		url := mainHtmlArticle.Request.URL.String()
		title := mainHtmlArticle.ChildText("h1")
		content, _ := mainHtmlArticle.DOM.Html()

		page.Id = url[len(url)-9:]
		page.Url = url
		page.Title = title
		page.Content = content
		book = append(book, page)

		fmt.Println("-------")
		fmt.Println("Id:", page.Id)
		fmt.Println("Link:", page.Url)
		fmt.Println("Title:", page.Title)
		//fmt.Println("Content:", item.Content)

		converter := md.NewConverter("", true, nil)
		converter.Use(plugin.GitHubFlavored())
		markdown, err := converter.ConvertString(page.Content)
		if err != nil {
			fmt.Println(err)
			os.Exit(10)
		}

		file, err := os.Create("./website/docs/kbs/" + page.Id + ".md")
		if err != nil {
			fmt.Println(err)
			os.Exit(10)
		} else {
			file.WriteString(markdown)
		}
		file.Close()
	})

	titleCollector.OnRequest(func(response *colly.Request) {
		fmt.Println("Visiting", response.URL.String())
	})
	articleCollector.OnRequest(func(response *colly.Request) {
		fmt.Println("Visiting Article", response.URL.String())
	})
	titleCollector.OnResponse(func(response *colly.Response) {
		fmt.Println("Got a response from", response.Request.URL)
	})

	titleCollector.OnHTML(".results_summary a[href]", func(nextPageHtmlLink *colly.HTMLElement) {
		nextPageUrl := nextPageHtmlLink.Request.AbsoluteURL(nextPageHtmlLink.Attr("href"))
		titleCollector.Visit(nextPageUrl)
	})

	titleCollector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Got this error:", err)
	})
	articleCollector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Got this error:", err)
	})

	titleCollector.OnScraped(func(response *colly.Response) {
		fmt.Println("Finished", response.Request.URL)
		js, err := json.MarshalIndent(book, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("book.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	titleCollector.Visit("https://www.suse.com/support/kb/?id=SUSE+Rancher")
}
