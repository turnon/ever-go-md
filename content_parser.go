package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type contentParser interface {
	ContentString() string
	Excerpt() string
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
				e.addParagraph(&imgRef{post: e.post, img: node})
			}

			return
		}

		e.addParagraph(&text{div})
	})
}

func (e *extracter) Excerpt() string {
	return ""
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

func (r *replacer) Excerpt() string {
	divs := r.RawBody().Find("div")
	text := strings.ReplaceAll(divs.First().Text(), `"`, `\"`)
	runeStr := []rune(text)

	if len(runeStr) > 140 {
		runeStr = runeStr[:140]
	}

	return string(runeStr) + "..."
}

func (r *replacer) ContentString() string {
	rawBody := r.RawBody()

	rawBody.Find("div[style]").Each(func(i int, div *goquery.Selection) {
		if style, _ := div.Attr("style"); strings.Index(style, "break-word") < 0 {
			return
		}

		div.BeforeSelection(div.Children())
		div.Remove()
	})

	rawBody.Find("img").Each(func(i int, img *goquery.Selection) {
		path := (&imgRef{post: r.post, img: img}).path()
		img.SetAttr("src", path)
	})

	rawBody.Find("a > img").Each(func(i int, img *goquery.Selection) {
		a := img.Parent()
		path := (&attachmentRef{r.post, a}).path()
		a.SetAttr("href", path)
	})

	rawBody.Find("div[style]").Each(func(i int, div *goquery.Selection) {
		if style, _ := div.Attr("style"); strings.Index(style, "box-sizing") < 0 {
			return
		}

		lines := []string{}
		div.Children().Each(func(i int, line *goquery.Selection) {
			lines = append(lines, line.Text())
		})
		code := strings.Join(lines, "\n")

		pre := strings.Join([]string{"<pre><code>", code, "</code></pre>"}, "")
		div.AfterHtml(pre)
		div.Remove()
	})

	html, err := rawBody.Html()
	if err != nil {
		panic(err)
	}
	return html
}
