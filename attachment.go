package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const assetsFiles = "/files"

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

	destPath := filepath.Join(dir, renamedAttachment(a.name()))
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
	*post
	a *goquery.Selection
}

func (a *attachmentRef) path() string {
	href, exists := a.a.Attr("href")
	if !exists {
		err := errors.New("href not found")
		panic(err)
	}

	newHref := strings.Replace(href, a.originAttachmentsSubDir(), a.assetsLocation(), 1)
	basename := filepath.Base(newHref)
	newBasename := renamedAttachment(basename)

	return strings.Replace(newHref, basename, newBasename, 1)
}

type imgRef struct {
	*post
	img       *goquery.Selection
	_fileName string
}

func (i *imgRef) fileName() string {
	if i._fileName != "" {
		return i._fileName
	}

	src, exists := i.img.Attr("src")
	if !exists {
		err := errors.New("src not found")
		panic(err)
	}

	i._fileName = renamedAttachment(filepath.Base(src))
	return i._fileName
}

func (i *imgRef) String() string {
	return `![alt text](` + i.path() + ` "` + i.fileName() + `")`
}

func (i *imgRef) path() string {
	return filepath.Join(i.assetsLocation(), i.fileName())
}

func renamedAttachment(path string) string {
	ext := filepath.Ext(path)
	noExt := strings.TrimSuffix(path, ext)
	return md5Str(noExt) + ext
}

func md5Str(str string) string {
	sum := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}
