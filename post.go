package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	p := &post{path: path, html: determinedFormat(data)}
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
	return p.baseName() + ".md"
}

func (p *post) baseName() string {
	name := filepath.Base(p.path)
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func (p *post) dirName() string {
	return filepath.Dir(p.path)
}

func (p *post) originAttachmentsDir() string {
	return filepath.Join(p.dirName(), p.originAttachmentsSubDir())
}

func (p *post) originAttachmentsSubDir() string {
	return p.baseName() + "_files"
}

func (p *post) copyAttachmentsTo(destDir string) {
	srcDir := p.originAttachmentsDir()
	if _, err := os.Stat(srcDir); err != nil && os.IsNotExist(err) {
		return
	}

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}

	destDir = filepath.Join(destDir, p.slug())
	if err := os.Mkdir(destDir, 0644); err != nil {
		panic(err)
	}

	for _, file := range files {
		a := attachment{path: filepath.Join(srcDir, file.Name())}
		a.copyToDir(destDir)
	}
}

func (p *post) slug() string {
	sum := md5.Sum([]byte(p.Title()))
	return fmt.Sprintf("%x", sum)
}

func (p *post) meta() string {
	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString("title: \"")
	sb.WriteString(p.Title())
	sb.WriteString("\"\n")
	sb.WriteString("slug: \"")
	sb.WriteString(p.slug())
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

		if nodeName == "br" {
			p.addParagraph(&br{})
			return
		}

		innerText := div.Text()

		if len(innerText) == 0 {
			if nodeName == "a" {
				p.addParagraph(&attachmentRef{p, node})
			}

			if nodeName == "img" {
				p.addParagraph(&imgRef{p, node})
			}

			return
		}

		p.addParagraph(&text{div})
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

	return strings.Join([]string{p.meta(), p.content()}, "\n\n")
}

func (p *post) content() string {
	strs := []string{}
	for _, p := range p.paragraphs {
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

func (l *link) String() string {
	href, _ := l.div.Children().First().Attr("href")
	return href
}
