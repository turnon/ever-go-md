package main

import (
	"bytes"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type html interface {
	Title() string
	Body() *goquery.Selection
}

type winHTML struct {
	data []byte
}

type macHTML struct {
	data []byte
}

var (
	winHTMLRebundantTags = regexp.MustCompile(`(?s)<a name="\d+"/>.*?<br/>`)
	macTitle             = regexp.MustCompile(`(?s)<title>(.*?)</title>`)
)

func (w *winHTML) Body() *goquery.Selection {
	data := winHTMLRebundantTags.ReplaceAll(w.data, []byte(""))
	doc := fileToDoc(data)
	return doc.Find("div > span").Children()
}

func (w *winHTML) Title() string {
	return ""
}

func (m *macHTML) Body() *goquery.Selection {
	doc := fileToDoc(m.data)
	return doc.Find("body").Children()
}

func (m *macHTML) Title() string {
	byts := macTitle.FindSubmatch(m.data)[1]
	return string(byts)
}

func fileToDoc(data []byte) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return doc
}
