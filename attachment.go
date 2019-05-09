package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type attachment struct {
	path string
}

func (a *attachment) copyToDir(dir string) {
	var err error
	var srcfd *os.File
	var destfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(a.path); err != nil {
		panic(err)
	}
	defer srcfd.Close()

	destPath := filepath.Join(dir, a.name())
	if destfd, err = os.Create(destPath); err != nil {
		panic(err)
	}
	defer destfd.Close()

	if _, err := io.Copy(destfd, srcfd); err != nil {
		panic(err)
	}

	if srcinfo, err = os.Stat(a.path); err != nil {
		panic(err)
	}

	if err := os.Chmod(destPath, srcinfo.Mode()); err != nil {
		panic(err)
	}
}

func (a *attachment) name() string {
	return filepath.Base(a.path)
}

type attachmentRef struct {
	subDir string
	div    *goquery.Selection
}

func (a *attachmentRef) String() string {
	return ""
}

type imgRef struct {
	subDir string
	img    *goquery.Selection
}

func (i *imgRef) String() string {
	filename, exists := i.img.Attr("data-filename")
	if !exists {
		panic("data-filename not found")
	}
	return `![alt text](/attachments/` + i.subDir + `/` + filename + ` "` + filename + `")`
}
