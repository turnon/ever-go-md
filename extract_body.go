package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func extractBody(p *post) string {
	parseBody := func(p *post) []paragraph {
		var paragraphs []paragraph

		p.Body().Each(func(i int, div *goquery.Selection) {
			if _, exists := div.Attr("style"); exists {
				paragraphs = append(paragraphs, &code{div})
				return
			}

			node := div.Children().First()
			nodeName := goquery.NodeName(node)

			if nodeName == "br" {
				paragraphs = append(paragraphs, &br{})
				return
			}

			innerText := div.Text()

			if len(innerText) == 0 {
				if nodeName == "a" {
					paragraphs = append(paragraphs, &attachmentRef{p, node})
				}

				if nodeName == "img" {
					paragraphs = append(paragraphs, &imgRef{p, node})
				}

				return
			}

			paragraphs = append(paragraphs, &text{div})
		})

		return paragraphs
	}

	strs := []string{}
	for _, p := range parseBody(p) {
		str := p.String()
		length := len(strs)
		if length > 0 && strs[length-1] == "\n" && str == "\n" {
			continue
		}

		strs = append(strs, str)
	}
	return strings.Join(strs, "\n\n")
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
