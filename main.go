package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	singleFile := flag.String("file", "", "the single file will be converted")
	fromDir := flag.String("from", "", "files in this dir will be converted")
	toDir := flag.String("to", "", "output to content/posts/ and static/attachments/ this dir if specified, else output to stdout")
	parserName := flag.String("parser", "replaceBody", "parser name, `extractBody` by default")
	clean := flag.Bool("clean", false, "clean destination dir")
	help := flag.Bool("help", false, "print usage")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	outFunc := output(*toDir, *clean)
	parse := postParser(*parserName)

	// handle single file
	if singleFilePath := *singleFile; singleFilePath != "" {
		p := parse(singleFilePath)
		outFunc(p)
		return
	}

	// handle multiple files
	formDirName := *fromDir
	files, err := ioutil.ReadDir(formDirName)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		p := parse(filepath.Join(formDirName, file.Name()))
		outFunc(p)
	}
}

func output(toDir string, clean bool) func(p *post) {
	if toDir == "" {
		return func(p *post) {
			fmt.Println(p.String())
		}
	}

	postsDir, attachmentsDir := filepath.Join(toDir, "_posts"), filepath.Join(toDir, assetsFiles)
	cleanDir(clean, postsDir)
	cleanDir(clean, attachmentsDir)

	return func(p *post) {
		dest := filepath.Join(postsDir, p.MdFileName())
		if err := ioutil.WriteFile(dest, []byte(p.String()), 0644); err != nil {
			panic(err)
		}

		p.copyAttachmentsTo(attachmentsDir)
	}
}

func cleanDir(clean bool, dir string) {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0644); err != nil {
			panic(err)
		}
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	if len(files) <= 0 {
		return
	}

	if !clean {
		panic(dir + " is not empty !")
	}

	for _, file := range files {
		err := os.RemoveAll(filepath.Join(dir, file.Name()))
		if err != nil {
			panic(err)
		}
	}
}
