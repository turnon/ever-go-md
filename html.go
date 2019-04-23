package main

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type html interface {
	Title() string
	CreatedAt() string
	Tags() []string
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
	title                = regexp.MustCompile(`(?s)<title>(.*?)</title>`)
	winCreatedAt         = regexp.MustCompile(`(?s)<tr><td><b>创建时间：</b></td><td><i>(.*?)\s.*?</i></td></tr>`)
	macCreatedAt         = regexp.MustCompile(`(?s)<meta name="created" content="(\d{4}-\d{2}-\d{2}).*?/>`)
	winTags              = regexp.MustCompile(`<tr><td><b>标签：</b></td><td><i>(.*?)</i></td></tr>`)
	macTags              = regexp.MustCompile(`(?s)<meta name="keywords" content="(.*?)"/>`)
)

func (w *winHTML) Body() *goquery.Selection {
	data := winHTMLRebundantTags.ReplaceAll(w.data, []byte(""))
	doc := fileToDoc(data)
	return doc.Find("div > span").Children()
}

func (w *winHTML) Title() string {
	byts := title.FindSubmatch(w.data)[1]
	return string(byts)
}

func (w *winHTML) CreatedAt() string {
	byts := winCreatedAt.FindSubmatch(w.data)[1]
	return string(byts)
}

func (w *winHTML) Tags() []string {
	byts := winTags.FindSubmatch(w.data)[1]
	return strings.Split(string(byts), ", ")
}

func (m *macHTML) Body() *goquery.Selection {
	doc := fileToDoc(m.data)
	return doc.Find("body").Children()
}

func (m *macHTML) Title() string {
	byts := title.FindSubmatch(m.data)[1]
	return string(byts)
}

func (m *macHTML) CreatedAt() string {
	byts := macCreatedAt.FindSubmatch(m.data)[1]
	return string(byts)
}

func (m *macHTML) Tags() []string {
	byts := macTags.FindSubmatch(m.data)[1]
	return strings.Split(string(byts), ", ")
}

func fileToDoc(data []byte) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return doc
}
