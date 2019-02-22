package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type (
	data []iconData

	iconData struct {
		BaseFileName string `json:"baseFileName"`
		Extension    string `json:"extension"`
		Sizes        []int  `json:"sizes"`
		Rel          string `json:"rel"`
	}
)

func main() {
	filePath := flag.String("filepath", "", "Path of the base image file")
	destPath := flag.String("dest", "", "Folder path of the destination")
	flag.Parse()
	if *filePath == "" || *destPath == "" {
		log.Fatal("enter the parameters")
	}
	src, err := imaging.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	jsonFile, err := os.Open("cmd/data.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	var d data
	b, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(b, &d)

	// 	_ = ioutil.WriteFile("", nil, os.ModeExclusive)

	var links []string

	for _, item := range d {
		for _, size := range item.Sizes {
			f := imaging.Resize(src, size, size, imaging.Lanczos)
			filename := fmt.Sprintf("%s-%dx%d.%s", item.BaseFileName, size, size, item.Extension)
			err = imaging.Save(f, path.Join(*destPath, filename))
			if err != nil {
				log.Fatal(err)
			}
			links = append(links, fmt.Sprintf(`{rel: "%s", sizes: "%dx%d", href: "/icons/%s"}`, item.Rel, size, size, filename))
		}
	}
	content := strings.Join(links, ",\n")
	of, err := os.Create("cmd/links.json")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()

	of.WriteString("[")
	of.WriteString(content)
	of.WriteString("]")
}
