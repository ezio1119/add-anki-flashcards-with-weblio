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
	Examples    examples
	AudioURL    string
	Level       int
}

type example struct {
	En, Ja string
}

type examples []*example

func (exx examples) String() string {
	var result string
	for i, ex := range exx {
		result += fmt.Sprintf("<strong>%s</strong>: %s", ex.En, ex.Ja)

		if len(exx) != i+1 {
			result += "<br>"
		}
	}

	return result
}

func Query(ctx context.Context, query string) (*QueryResult, error) {
	r := &QueryResult{
		Query: query,
	}
	url := fmt.Sprintf("https://ejje.weblio.jp/content/%s", r.Query)

	var doc *html.Node
	docChan := make(chan *html.Node)
	errChan := make(chan error)

	go func() {
		doc, err := htmlquery.LoadURL(url)
		if err != nil {
			errChan <- err
			return
		}
		docChan <- doc
	}()

	select {
	case doc = <-docChan:

	case err := <-errChan:
		return nil, fmt.Errorf("weblio: Query: failed load url: %s: %w", r.Query, err)

	case <-ctx.Done():
		return nil, fmt.Errorf("weblio: Query: failed load url: %s: %w", r.Query, ctx.Err())
	}

	if r.Description = getDescription(doc); r.Description == "" {
		return nil, fmt.Errorf("weblio: Query: failed get description: %s", r.Query)
	}

	if r.AudioURL = getAudioURL(doc); r.AudioURL == "" {
		fmt.Printf("weblio: Query: failed get audioURL: %s", r.Query)
	}

	var err error
	r.Level, err = getLevel(doc)
	if err != nil || r.Level == 0 {
		fmt.Printf("weblio: Query: failed get level: %s: %s", r.Query, err)
	}

	if r.Examples = getExamples(doc); len(r.Examples) == 0 {
		fmt.Printf("weblio: Query: failed get example: %s", r.Query)
	}

	return r, nil
}

func getQuery(doc *html.Node) string {
	queryNode := htmlquery.FindOne(doc, "//*[@id=\"h1Query\"]")
	if queryNode != nil {
		return htmlquery.InnerText(queryNode)
	}

	return ""
}

func getDescription(doc *html.Node) (desc string) {
	defer func() {
		if desc != "" {
			desc = strings.TrimSpace(desc)
		}
	}()

	descNode := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p/span[2]")
	if descNode != nil {
		desc = htmlquery.InnerText(descNode)
	}

	descNode = htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p[2]/span")
	if descNode != nil {
		desc = htmlquery.InnerText(descNode)
	}

	return desc
}

func getAudioURL(doc *html.Node) string {
	audioNode := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/table[1]/tbody/tr/td[2]/table/tbody/tr[1]/td[1]/i/audio/source")
	if audioNode != nil {
		return htmlquery.SelectAttr(audioNode, "src")
	}

	return ""
}

func getLevel(doc *html.Node) (int, error) {
	levelNode := htmlquery.FindOne(doc, "//*[@id=\"learning-level-table\"]/div/span[1]/span[3]")
	if levelNode == nil {
		return 0, nil
	}

	levelStr := htmlquery.InnerText(levelNode)
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		return 0, err
	}

	return level, nil
}

func getExamples(doc *html.Node) examples {
	exx := examples{}

	exampleNodes := htmlquery.Find(doc, "//*[contains(concat(' ',@class,' '),' KejjeYrLn ')]")

	for _, e := range exampleNodes {
		var en, ja string

		enNode := e.FirstChild
		jaNode := e.LastChild

		for _, a := range enNode.Attr {
			if a.Val != "KejjeYrEn" {
				continue
			}

			for wordNode := enNode.FirstChild; wordNode != nil; wordNode = wordNode.NextSibling {
				en += htmlquery.InnerText(wordNode)
			}
		}

		for _, a := range jaNode.Attr {
			if a.Val != "KejjeYrJp" {
				continue
			}

			for wordNode := jaNode.FirstChild; wordNode != nil; wordNode = wordNode.NextSibling {
				ja += htmlquery.InnerText(wordNode)
			}
		}

		if en == "" || ja == "" {
			continue
		}

		exx = append(exx, &example{en, ja})
	}

	return exx
}
