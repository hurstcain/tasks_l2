package main

import (
	"github.com/opesun/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Wget struct {
	path string
}

func NewWget(path string) *Wget {
	return &Wget{
		path: path,
	}
}

func (w Wget) GetHtml(url, path string) {
	html, err := goquery.ParseUrl(url)
	if err != nil {
		log.Printf("Error when parsing url: %s\n", err.Error())
		return
	}
	fileName := path + "index.html"
	w.WriteToFile(fileName, html.Html())
	w.ResourcesSearch(html)
}

func (w Wget) WriteToFile(fileName, html string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error when opening file: %s\n", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		log.Printf("Error when writing to file: %s\n", err.Error())
		return
	}
}

func (w Wget) ResourcesSearch(html goquery.Nodes) {
	path := w.path + `resources\`
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Printf("Error when cteating directory: %s\n", err.Error())
		return
	}

	for _, src := range html.Find("").Attrs("src") {
		urlParts := strings.Split(src, "/")
		filename := path + urlParts[len(urlParts)-1]
		w.Download(filename, src)
	}

	for _, href := range html.Find("").Attrs("href") {
		urlParts := strings.Split(href, "/")
		filename := path + urlParts[len(urlParts)-1]
		w.Download(filename, href)
	}
}

func (w Wget) Download(fileName, url string) {
	resource, err := http.Get(url)
	if err != nil {
		return
	}

	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resource.Body)
	if err != nil {
		return
	}

	log.Printf("Resource %s downloaded\n", url)
}

func main() {
	argsLen := len(os.Args)
	url := os.Args[argsLen-1]
	path, err := os.Getwd()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	path = path + `\site\`
	os.Mkdir(path, 0755)

	wget := NewWget(path)
	wget.GetHtml(url, path)
}
