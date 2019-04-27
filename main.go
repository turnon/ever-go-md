package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fromDir, toDir := os.Args[1], os.Args[2]
	files, err := ioutil.ReadDir(fromDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := filepath.Join(fromDir, file.Name())
		data, err := ioutil.ReadFile(src)
		if err != nil {
			panic(err)
		}

		p := newPost(data)

		dest := filepath.Join(toDir, file.Name())
		dest = strings.Replace(dest, ".html", ".md", 1)
		if err := ioutil.WriteFile(dest, []byte(p.String()), file.Mode()); err != nil {
			panic(err)
		}

	}

}
