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

func (p *post) parseBody(divs *goquery.Selection) {
	divs.Each(func(i int, div *goquery.Selection) {
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
	})
}

func (p *post) addParagraph(para paragraph) {
	p.paragraphs = append(p.paragraphs, para)
}

func (p *post) String() string {
	strs := []string{}
	for _, p := range p.paragraphs {
		str := p.String()
		length := len(strs)
		if length > 0 && strs[length-1] == "\n" && str == "\n" {
			continue
		}

		strs = append(strs, str)
	}
	return strings.Join(strs, "\n")
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

	sb.WriteString("```")
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
