package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// fmt.Println("Getting...")
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
			fmt.Println(k, v)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// c.Wait()
	c.Visit("https://wallhaven.cc/random")
}
