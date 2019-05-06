package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type post struct {
	path string
	html
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

func newPostFromPath(path string) *post {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	p := newPost(data)
	p.path = path
	return p
}

func newPost(data []byte) *post {
	p := &post{html: determinedFormat(data)}
	p.parse()
	return p
}

func determinedFormat(data []byte) html {
	if runtime.GOOS == "windows" {
		return &winHTML{data}
	}

	return &macHTML{data}
}

func (p *post) parse() {
	p.parseBody()
}

func (p *post) MdFileName() string {
	sum := md5.Sum([]byte(p.String()))
	return fmt.Sprintf("%x", sum) + ".md"
}

func (p *post) meta() string {
	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString("title: \"")
	sb.WriteString(p.Title())
	sb.WriteString("\"\n")
	sb.WriteString("date: ")
	sb.WriteString(p.CreatedAt())
	sb.WriteString("\n")
	sb.WriteString("tags: ")
	sb.WriteString(p.tagsStr())
	sb.WriteString("\n")
	sb.WriteString("---\n")
	return sb.String()
}

func (p *post) tagsStr() string {
	return `["` + strings.Join(p.Tags(), `", "`) + `"]`
}

func (p *post) parseBody() {
	p.Body().Each(func(i int, div *goquery.Selection) {
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
	defer func() {
		if r := recover(); r != nil {
			r = errors.Wrap(r.(runtime.Error), p.path)
			panic(r)
		}
	}()

	strs := []string{p.meta()}
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
