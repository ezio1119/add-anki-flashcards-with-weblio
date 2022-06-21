package weblio

import (
	"context"
	"fmt"
	"strconv"

	"github.com/antchfx/htmlquery"
)

type QueryResult struct {
	Query       string
	Description string
	AudioURL    string
	Level       int
}

func Query(ctx context.Context, query string) (*QueryResult, error) {
	result := &QueryResult{}

	url := fmt.Sprintf("https://ejje.weblio.jp/content/%s", query)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	queryElem := htmlquery.FindOne(doc, "//*[@id=\"h1Query\"]")
	if queryElem != nil {
		result.Query = htmlquery.InnerText(queryElem)
	} else {
		return nil, fmt.Errorf("weblio: Query: failed get query: %s", query)
	}
	//
	descElem := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p/span[2]")
	if descElem != nil {
		result.Description = htmlquery.InnerText(descElem)
	} else {
		descElem = htmlquery.FindOne(doc, "//*[@id=\"summary\"]/div[2]/p[2]/span")
		if descElem != nil {
			result.Description = htmlquery.InnerText(descElem)
		} else {
			return nil, fmt.Errorf("weblio: Query: failed get description: %s", query)
		}
	}

	audioElem := htmlquery.FindOne(doc, "//*[@id=\"summary\"]/table[1]/tbody/tr/td[2]/table/tbody/tr[1]/td[1]/i/audio/source")
	if audioElem != nil {
		result.AudioURL = htmlquery.SelectAttr(audioElem, "src")
	}

	levelElem := htmlquery.FindOne(doc, "//*[@id=\"learning-level-table\"]/div/span[1]/span[3]")
	if levelElem != nil {
		levelStr := htmlquery.InnerText(levelElem)

		result.Level, err = strconv.Atoi(levelStr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
