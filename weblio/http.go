package weblio

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type QueryResult struct {
	Query       string
	Description string
	AudioURL    string
	Level       int
}

func Query(ctx context.Context, query string) (*QueryResult, error) {
	r := &QueryResult{
		Query: query,
	}

	url := fmt.Sprintf("https://ejje.weblio.jp/content/%s", r.Query)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	if r.Description = getDescription(doc); r.Description == "" {
		return nil, fmt.Errorf("weblio: Query: failed get description: %s", r.Query)
	}

	if r.AudioURL = getAudioURL(doc); r.AudioURL == "" {
		fmt.Printf("weblio: Query: failed get audioURL: %s", r.Query)
	}

	r.Level, err = getLevel(doc)
	if err != nil || r.Level == 0 {
		fmt.Printf("weblio: Query: failed get level: %s: %s", r.Query, err)
	}

	return r, nil
}

func getQuery(doc *html.Node) string {
	queryElem := htmlquery.FindOne(doc, "//*[@id=\"h1Query\"]")
	if queryElem != nil {
		return htmlquery.InnerText(queryElem)
	}

	return ""
}

func getDescription(doc *html.Node) (desc string) {
	defer func() {
		if desc != "" {
			fmt.Println(desc)
			desc = strings.TrimSpace(desc)
			fmt.Println(desc)
		}
	}()

	descElem := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p/span[2]")
	if descElem != nil {
		desc = htmlquery.InnerText(descElem)
	}

	descElem = htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p[2]/span")
	if descElem != nil {
		desc = htmlquery.InnerText(descElem)
	}

	return desc
}

func getAudioURL(doc *html.Node) string {
	audioElem := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/table[1]/tbody/tr/td[2]/table/tbody/tr[1]/td[1]/i/audio/source")
	if audioElem != nil {
		return htmlquery.SelectAttr(audioElem, "src")
	}

	return ""
}

func getLevel(doc *html.Node) (int, error) {
	levelElem := htmlquery.FindOne(doc, "//*[@id=\"learning-level-table\"]/div/span[1]/span[3]")
	if levelElem != nil {
		levelStr := htmlquery.InnerText(levelElem)

		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return 0, err
		}

		return level, nil
	}

	return 0, nil
}
