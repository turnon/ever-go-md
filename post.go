package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type post struct {
	path string
	htmlFile
	contentParser
}

func pickContentParser(p *post, contentParserName string) {
	if contentParserName == "extractBody" {
		p.contentParser = &extracter{post: p}
	} else {
		p.contentParser = &replacer{post: p}
	}
}

func postParser(contentParserName string) func(path string) *post {
	return func(path string) *post {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		p := &post{path: path, htmlFile: determinedFormat(data)}
		pickContentParser(p, contentParserName)
		return p
	}
}

func determinedFormat(data []byte) htmlFile {
	if runtime.GOOS == "windows" {
		return &winHTML{data}
	}

	return &macHTML{data}
}

func (p *post) MdFileName() string {
	return p.CreatedAt() + "-" + p.baseName() + ".md"
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
	metas := []string{
		"---\n",
		"title: \"", p.Title(), "\"\n",
		"slug: \"", p.slug(), "\"\n",
		"date: ", p.CreatedAt(), "\n",
		"excerpt: \"", p.Excerpt(), "\"\n",
		"tags: ", p.tagsStr(), "\n",
		"---\n",
	}
	return strings.Join(metas, "")
}

func (p *post) tagsStr() string {
	return `["` + strings.Join(p.Tags(), `", "`) + `"]`
}

func (p *post) String() string {
	defer func() {
		if r := recover(); r != nil {
			r = errors.Wrap(r.(error), p.path)
			panic(r)
		}
	}()

	return strings.Join([]string{p.meta(), p.ContentString()}, "\n\n")
}
