package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	file := os.Args[1]

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	p := newPost(data)

	fmt.Println(len(p.paragraphs))
	fmt.Println(p.String())
}
