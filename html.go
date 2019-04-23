package main

import (
	"bytes"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type html interface {
	Body() *goquery.Selection
}

type winHTML struct {
	data []byte
}

type macHTML struct {
	data []byte
}

var winHTMLRebundantTags = regexp.MustCompile(`(?s)<a name="\d+"/>.*?<br/>`)

func (w *winHTML) Body() *goquery.Selection {
	data := winHTMLRebundantTags.ReplaceAll(w.data, []byte(""))
	doc := fileToDoc(data)
	return doc.Find("div > span").Children()
}

func (m *macHTML) Body() *goquery.Selection {
	doc := fileToDoc(m.data)
	return doc.Find("body").Children()
}

func fileToDoc(data []byte) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return doc
}
