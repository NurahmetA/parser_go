package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

type JsonData struct {
	Id       int    `json:"id"`
	Date     string `json:"date"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Author   string `json:"author"`
	Code     string `json:"code"`
}

func Parse() ([]JsonData, error) {
	var jsonData []JsonData
	var err error
	htmlRes, err := getHtmlPage()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlRes))
	doc.Find("#exploits-table tbody tr").Each(func(j int, item *goquery.Selection) {
		var data JsonData
		item.Children().Each(func(j int, child *goquery.Selection) {
			if j == 0 {
				data.Date = child.Text()
			}
			if j == 4 {
				data.Title = child.Children().Text()
				codeLink, _ := child.Children().Attr("href")
				data.Code = "https://exploit-db.com" + codeLink
			}
			if j == 5 {
				data.Type = child.Children().Text()
			}
			if j == 6 {
				data.Platform = child.Children().Text()
			}
			if j == 7 {
				data.Author = child.Children().Text()
			}
		})
		jsonData = append(jsonData, data)
	})
	return jsonData, err
}

func getHtmlPage() (string, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	var htmlRes string
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://exploit-db.com`),
		chromedp.Sleep(time.Second),
		chromedp.Click("verifiedCheck", chromedp.ByID),
		chromedp.Sleep(time.Second),
		chromedp.WaitVisible("tbody tr", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			htmlRes, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	return htmlRes, err
}
