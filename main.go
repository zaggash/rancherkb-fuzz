package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	mdPlugin "github.com/JohannesKaufmann/html-to-markdown/plugin"
	colly "github.com/gocolly/colly/v2"
	collyQueue "github.com/gocolly/colly/v2/queue"
	sf "github.com/simpleforce/simpleforce"
)

var (
	sfURL        = os.Getenv("sfURL")
	sfAPIVersion = "56.0"
	sfUser       = os.Getenv("sfUser")
	sfPassword   = os.Getenv("sfPass")
	sfToken      = os.Getenv("sfToken")
	sfIdsRequest = os.Getenv("sfIdsRequest")
	kbPath       = "./website/docs/kbs/"
	logger       *log.Logger
)

// Article skeleton
type Article struct {
	Id      string
	Title   string
	Url     string
	Content string
}

// Return simpleforce authenticated client
func createSfClient() *sf.Client {
	client := sf.NewClient(sfURL, sf.DefaultClientID, sfAPIVersion)
	if client == nil {
		logger.Fatal("error creating salesforce client")
		return nil
	}

	err := client.LoginPassword(sfUser, sfPassword, sfToken)
	if err != nil {
		logger.Fatal("failed with salesforce client authentication")
		return nil
	}
	return client
}

// Return SOQL Query results
func sfQuery(q string) *sf.QueryResult {
	sfClient := createSfClient()

	result, err := sfClient.Query(q)
	if err != nil {
		logger.Fatal("query failed")
	}

	if result.TotalSize < 1 {
		logger.Println("no records returned.")
	}

	return result
}

// Do Main stuff
func main() {
	// setup scrapper logger
	logger = log.New(os.Stderr, "[RancherKB-Fuzz] ", log.Lmsgprefix|log.LstdFlags)

	// get articles ID
	articlesId := sfQuery(sfIdsRequest)

	// setup articles list
	book := make([]Article, 0)

	// setup GoColly Queue
	collyQ, _ := collyQueue.New(
		10,
		&collyQueue.InMemoryQueueStorage{MaxSize: 10000},
	)

	// prepare articles Collector
	articleCollector := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`www.suse.com/support/kb/(.*)/?id(.*)`),
		),
		colly.MaxDepth(1),
	)

	articleCollector.Limit(&colly.LimitRule{
		Parallelism: 4,
		RandomDelay: 5 * time.Second,
	})

	articleCollector.WithTransport(&http.Transport{
		DisableKeepAlives: false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          50,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
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
		fmt.Println("-------")

		converter := md.NewConverter("", true, nil)
		converter.Use(mdPlugin.GitHubFlavored())
		markdown, err := converter.ConvertString(page.Content)
		if err != nil {
			logger.Println(err)
			os.Exit(5)
		}

		file, err := os.Create(kbPath + page.Id + ".md")
		if err != nil {
			logger.Println(err)
			os.Exit(10)
		} else {
			file.WriteString(markdown)
		}
		file.Close()
	})

	articleCollector.OnRequest(func(response *colly.Request) {
		logger.Println("visiting kb:", response.URL.String())
	})

	// articleCollector.OnResponse(func(response *colly.Response) {
	// 	fmt.Println("Response received:", response.StatusCode)
	// })

	articleCollector.OnError(func(response *colly.Response, err error) {
		logger.Println("Http Code:", response.StatusCode)
		logger.Println("got this error:", err)
		os.Exit(15)
	})

	for _, record := range articlesId.Records {
		id := record.StringField("ArticleNumber")
		collyQ.AddURL("https://www.suse.com/support/kb/doc/?id=" + id)
	}

	collyQ.Run(articleCollector)
}
