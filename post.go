package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type post struct {
	paragraphs []paragraph
}

type paragraph interface {
	String() string
}

type code struct {
	div *goquery.Selection
}

type text struct {
	div *goquery.Selection
}

type br struct {
}

type link struct {
	div *goquery.Selection
}

func (p *post) addParagraph(para paragraph) {
	p.paragraphs = append(p.paragraphs, para)
}

func (p *post) String() string {
	sb := strings.Builder{}
	for _, p := range p.paragraphs {
		sb.WriteString(p.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c *code) String() string {
	sb := strings.Builder{}
	sb.WriteString("```\n")

	c.div.Find("div").Each(func(i int, span *goquery.Selection) {
		if text := span.Text(); text != "" {
			sb.WriteString(text)
		}

		sb.WriteString("\n")
	})

	sb.WriteString("```\n")
	return sb.String()
}

func (t *text) String() string {
	return t.div.Text()
}

func (s *br) String() string {
	return "\n"
}

func (l *link) String() string {
	href, _ := l.div.Children().First().Attr("href")
	return href
}
