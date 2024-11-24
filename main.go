package main

import (
	"fmt"
	"log"
	"os"
	str "strings"
)

func usage() {
	fmt.Printf("Usage: %s inputdir [outputdir]\n", os.Args[0])
}

func filepathToString(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	stats, errrr := file.Stat()
	if errrr != nil {
		log.Fatal(errrr)
	}
	filesize := stats.Size()

	buf := make([]byte, filesize)
	_, errr := file.Read(buf)
	if errr != nil {
		log.Fatal(errr)
	}
	var builder str.Builder
	for _, b := range buf {
		builder.WriteByte(b)
	}
	return builder.String()
}

func markdownToHtml(md string) string {
	var builder str.Builder
	lines := str.Split(md, "\n")
	for _, line := range lines {
		if str.HasPrefix(line, "#") {
			hashcount := str.Count(line, "#")
			title := line[hashcount+1 : len(line)]
			builder.WriteString(fmt.Sprintf("<h%d>%s</h%d>", hashcount, title, hashcount))
			builder.WriteByte('\n')
		} else {
			builder.WriteString(line)
			builder.WriteByte('\n')
		}
	}
	s := builder.String()
	return s[:len(s)-1]
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	inputDir := os.Args[1]

	inputs, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	filePaths := []string{}
	for _, file := range inputs {
		if str.HasSuffix(file.Name(), ".md") {
			filePaths = append(filePaths, file.Name())
		}
	}
	for _, path := range filePaths {
		var builder str.Builder
		builder.WriteString(inputDir)
		builder.WriteString("/")
		builder.WriteString(path)
		filepath := builder.String()
		data := filepathToString(filepath)
		html := markdownToHtml(data)
		fmt.Printf("%s", html)
	}
}
