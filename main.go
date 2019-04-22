package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	file := os.Args[1]

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	p := post{}

	doc.Find("body").Children().Each(func(i int, div *goquery.Selection) {
		if _, exists := div.Attr("style"); exists {
			p.addParagraph(&code{div})
			return
		}

		node := div.Children().First()
		nodeName := goquery.NodeName(node)

		if nodeName == "" {
			p.addParagraph(&text{div})
			return
		}

		if nodeName == "br" {
			p.addParagraph(&br{})
			return
		}

		if nodeName == "a" {
			p.addParagraph(&link{div})
			return
		}

		fmt.Println(nodeName, node.Text())
	})

	fmt.Println(len(p.paragraphs))
	fmt.Println(p.String())
}
