package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

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

	p := &post{}

	if runtime.GOOS == "windows" {
		p.parseBody(doc.Find("body > div").Children())
	} else {
		p.parseBody(doc.Find("body").Children())
	}

	fmt.Println(len(p.paragraphs))
	fmt.Println(p.String())
}
