package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	file := os.Args[1]

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var divs *goquery.Selection

	if runtime.GOOS == "windows" {
		divs = windowsBody(data)
	} else {
		divs = macBody(data)
	}

	p := &post{}
	p.parseBody(divs)

	fmt.Println(len(p.paragraphs))
	fmt.Println(p.String())
}

func windowsBody(data []byte) *goquery.Selection {
	r, _ := regexp.Compile(`(?s)<a name="\d+"/>.*?<br/>`)
	data = r.ReplaceAll(data, []byte(""))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return doc.Find("div > span").Children()
}

func macBody(data []byte) *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return doc.Find("body").Children()
}
