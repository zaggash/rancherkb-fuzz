package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	sf "github.com/simpleforce/simpleforce"
)

var (
	sfAPIVersion = "62.0"
	sfURL        = os.Getenv("sfURL")
	sfUser       = os.Getenv("sfUser")
	sfPassword   = os.Getenv("sfPass")
	sfToken      = os.Getenv("sfToken")
	sfIdsRequest = os.Getenv("sfIdsRequest")

	kbPath   = "./website/docs/kbs/"
	logger   = log.New(os.Stderr, "[RancherKB-Fuzz] ", log.LstdFlags)
	sfClient *sf.Client
)

type Article struct {
	Id                string
	Title             string
	UrlName           string
	Content_Env2      string
	ContentBody       string
	ContentSituation  string
	ContentResolution string
	ContentCause      string
	ContentState      string
	ContentProducts   string
}

func init() {
	logger = log.New(os.Stderr, "[RancherKB-Fuzz] ", log.Lmsgprefix|log.LstdFlags)
}

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

func sfQuery(q string) *sf.QueryResult {
	if sfClient == nil {
		sfClient = createSfClient()
	}
	result, err := sfClient.Query(q)
	if err != nil {
		logger.Fatal("query failed")
	}
	if result.TotalSize < 1 {
		logger.Println("no records returned.")
	}

	allRecords := result.Records

	// Loop to fetch remaining records if the first page isn't the last one
	for !result.Done {
		result, err = sfClient.Query(result.NextRecordsURL)
		if err != nil {
			logger.Fatal("query more failed")
		}
		allRecords = append(allRecords, result.Records...)
	}
	result.Records = allRecords
	return result
}

func main() {
	articles := sfQuery(sfIdsRequest)

	for _, record := range articles.Records {
		article := Article{
			Id:                record.StringField("ArticleNumber"),
			Title:             record.StringField("Title"),
			UrlName:           record.StringField("UrlName"),
			Content_Env2:      record.StringField("SUSE_Environment2__c"),
			ContentBody:       record.StringField("Body__c"),
			ContentSituation:  record.StringField("SUSE_Situation__c"),
			ContentResolution: record.StringField("SUSE_Resolution__c"),
			ContentCause:      record.StringField("SUSE_Cause__c"),
			ContentState:      record.StringField("SUSE_State__c"),
			//ContentProducts:   record.StringField("SUSE_Products__c"),
		}

		contents := map[string]string{
			"Environment": article.Content_Env2,
			"Situation":   article.ContentSituation,
			"State":       article.ContentState,
			"Procedure":   article.ContentBody,
			"Cause":       article.ContentCause,
			"Resolution":  article.ContentResolution,
		}

		fmt.Println("-------")
		fmt.Println("Id:", article.Id)
		fmt.Println("Title:", article.Title)
		fmt.Println("UrlName: https://support.scc.suse.com/s/kb/" + article.UrlName)

		markdowns := make(map[string]string)
		for key, val := range contents {
			if val == "" {
				logger.Printf("Warning: section '%s' is empty, skipping conversion", key)
				continue
			}
			mdStr, err := md.ConvertString(val)
			if err != nil {
				logger.Println("error converting HTML to Markdown:", err)
				os.Exit(5)
			}
			markdowns[key] = mdStr
		}

		kbContent := fmt.Sprintf("# %s\n\n", article.Title)
		kbContent += fmt.Sprintf("**Article Number:** [%s](https://support.scc.suse.com/s/kb/%s)\n\n", article.Id, article.UrlName)
		// Add sections in a specific order, skipping any that are missing
		for _, section := range []string{"Environment", "Situation", "State", "Procedure", "Cause", "Resolution"} {
			if _, exists := markdowns[section]; !exists {
				logger.Printf("Warning: section '%s' is missing and will be skipped in the output", section)
				continue
			}
			kbContent += fmt.Sprintf("## **%s**\n\n%s\n\n", section, markdowns[section])
		}

		filePath := filepath.Join(kbPath, article.Id+".md")
		file, err := os.Create(filePath)
		if err != nil {
			logger.Println("error creating file:", err)
			return
		}
		_, err = file.WriteString(kbContent)
		if err != nil {
			logger.Println("error writing to file:", err)
			defer file.Close()
		}

		fmt.Println("-------")

	}

	files, err := os.ReadDir(kbPath)
	if err != nil {
		logger.Println("error reading directory:", err)
	} else {
		fmt.Printf("Total files created: %d\n", len(files))
	}
}
