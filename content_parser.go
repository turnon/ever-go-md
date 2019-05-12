package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type contentParser interface {
	ContentString() string
}

type extracter struct {
	*post
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

func (e *extracter) addParagraph(para paragraph) {
	e.paragraphs = append(e.paragraphs, para)
}

func (e *extracter) parseBody() {
	e.Body().Each(func(i int, div *goquery.Selection) {
		if _, exists := div.Attr("style"); exists {
			e.addParagraph(&code{div})
			return
		}

		node := div.Children().First()
		nodeName := goquery.NodeName(node)

		if nodeName == "br" {
			e.addParagraph(&br{})
			return
		}

		innerText := div.Text()

		if len(innerText) == 0 {
			if nodeName == "a" {
				e.addParagraph(&attachmentRef{e.post, node})
			}

			if nodeName == "img" {
				e.addParagraph(&imgRef{e.post, node})
			}

			return
		}

		e.addParagraph(&text{div})
	})
}

func (e *extracter) ContentString() string {
	e.parseBody()

	strs := []string{}
	for _, p := range e.paragraphs {
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

type replacer struct {
	*post
}

func (r *replacer) ContentString() string {
	html, err := r.RawBody().Html()
	if err != nil {
		panic(err)
	}
	return html
}