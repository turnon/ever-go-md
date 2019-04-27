package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	singleFile := flag.String("file", "", "file path")
	fromDir := flag.String("from", "", "from dir")
	toDir := flag.String("to", "", "to dir")
	flag.Parse()

	outFunc := output(*toDir)

	if singleFilePath := *singleFile; singleFilePath != "" {
		p := newPostFromPath(singleFilePath)
		outFunc(p)
		return
	}

	handleFiles(*fromDir, outFunc)
}

func handleFiles(fromDir string, outFunc func(p *post)) {
	files, err := ioutil.ReadDir(fromDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		p := newPostFromPath(filepath.Join(fromDir, file.Name()))
		outFunc(p)
	}
}

func output(toDir string) func(p *post) {
	if toDir == "" {
		return func(p *post) {
			fmt.Println(p.String())
		}
	}

	return func(p *post) {
		dest := filepath.Join(toDir, p.MdFileName())
		if err := ioutil.WriteFile(dest, []byte(p.String()), 0644); err != nil {
			panic(err)
		}
	}
}
