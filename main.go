package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/gocolly/colly"
)

var url = "https://wallhaven.cc/"

var args struct {
	Category string `help:"Category - random, latest, toplist"`
}

func main() {
	arg.MustParse(&args)
	url = url + args.Category
	scrape(url)

}

func scrape(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("wallhaven.cc"),
		colly.CacheDir("./__cache"),
	)
	links := make(map[string]string)

	c.OnHTML("figure.thumb", func(e *colly.HTMLElement) {
		link := e.ChildAttr("img", "data-src")
		re := regexp.MustCompile(`(\/th.)(wallhaven.cc)(\/small)(\/\w{2}\/)(\w{6}\.jpg)`)
		formatted := re.ReplaceAllString(link, `/w.$2/full${4}wallhaven-$5`)
		filename := strings.Split(formatted, "wallhaven-")[1]
		links[filename] = formatted
	})

	c.OnScraped(func(r *colly.Response) {
		for k, v := range links {
			fetchAndSave(k, v)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(url)
}

func fetchAndSave(filename string, url string) {
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		os.Mkdir("./images", os.ModePerm)
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create("./images/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File saved!", filename)
}
