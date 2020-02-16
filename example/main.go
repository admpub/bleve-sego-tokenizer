package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve"

	sego "github.com/admpub/bleve-sego-tokenizer"
)

func main() {
	// open a new index
	indexMapping := bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer("sego",
		map[string]interface{}{
			"files": "../dict.txt",
			"type":  sego.Name,
		})
	if err != nil {
		log.Fatal(err)
	}

	// create a custom analyzer
	err = indexMapping.AddCustomAnalyzer("sego",
		map[string]interface{}{
			"type":      sego.Name,
			"tokenizer": sego.Name,
		})

	if err != nil {
		log.Fatal(err)
	}

	indexMapping.DefaultAnalyzer = "sego"
	cacheDir := "sego.beleve"
	os.RemoveAll(cacheDir)
	index, err := bleve.New(cacheDir, indexMapping)

	if err != nil {
		log.Fatal(err)
	}

	docs := []struct {
		Title string
		Name  string
	}{
		{
			Title: "Doc 1",
			Name:  "This is the first document we’ve added",
		},
		{
			Title: "Doc 2",
			Name:  "The second one 你 中文测试中文 is even more interesting! 吃水果",
		},
		{
			Title: "Doc 3",
			Name:  "买水果然后来世博园。",
		},
		{
			Title: "Doc 4",
			Name:  "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作",
		},
		{
			Title: "Doc 5",
			Name:  "咱俩交换一下吧。",
		},
		{
			Title: "Doc 6",
			Name:  "富士康小时工周薪薪",
		},
		{
			Title: "Doc 7",
			Name:  "富康小时周薪",
		},
		{
			Title: "Doc 8",
			Name:  "富小周",
		},
	}
	// index docs
	for _, doc := range docs {
		index.Index(doc.Title, doc)
	}

	// search for some text
	for _, keyword := range []string{"水果世博园", "你", "first", "中文", "交换机", "交换", "富士"} {
		query := bleve.NewQueryStringQuery(keyword)
		search := bleve.NewSearchRequest(query)
		search.Highlight = bleve.NewHighlight()
		searchResults, err := index.Search(search)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result of \"%s\": %d matches:\n", keyword, searchResults.Total)
		for i, hit := range searchResults.Hits {
			rv := fmt.Sprintf("%d. %s, (%f)\n", i+searchResults.Request.From+1, hit.ID, hit.Score)
			for fragmentField, fragments := range hit.Fragments {
				rv += fmt.Sprintf("%s: ", fragmentField)
				for _, fragment := range fragments {
					rv += fmt.Sprintf("%s", fragment)
				}
			}
			fmt.Printf("%s\n", rv)
		}
		fmt.Printf("\n")
	}

	// Output:
	// Result of "水果世博园": 2 matches:
	// 1. Doc 3, (0.981169)
	// Name: 买<span class="highlight">水果</span>然后来<span class="highlight">世博</span>园。
	// 2. Doc 2, (0.051133)
	// Name: The second one 你 中文测试中文 is even more interesting! 吃<span class="highlight">水果</span>
	// Result of "你": 1 matches:
	// 1. Doc 2, (0.399574)
	// Name: The second one <span class="highlight">你</span> 中文测试中文 is even more interesting! 吃水果
	// Result of "first": 1 matches:
	// 1. Doc 1, (0.512150)
	// Name: This is the <span class="highlight">first</span> document we’ve added
	// Result of "中文": 1 matches:
	// 1. Doc 2, (0.565083)
	// Name: The second one 你 <span class="highlight">中文</span>测试<span class="highlight">中文</span> is even more interesting! 吃水果
	// Result of "交换机": 2 matches:
	// 1. Doc 4, (0.488048)
	// Name: 工信处女干事每月经过下属科室都要亲口交代24口<span class="highlight">交换机</span>等技术性器件的安装工作
	// 2. Doc 5, (0.165357)
	// Name: 咱俩<span class="highlight">交换</span>一下吧。
	// Result of "交换": 2 matches:
	// 1. Doc 5, (0.534158)
	// Name: 咱俩<span class="highlight">交换</span>一下吧。
	// 2. Doc 4, (0.302165)
	// Name: 工信处女干事每月经过下属科室都要亲口交代24口<span class="highlight">交换</span>机等技术性器件的安装工作
}
