package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = http.Transport{
		TLSClientConfig: config,
	}
	netClient = &http.Client{
		Transport: &transport,
	}

	queue = make(chan string)
)

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println(("Missing URL"))
		os.Exit(1)
	}

	go func() {
		queue <- arguments[0]
	}()
	//baseUrl := "https://exercism.org/tracks/go/exercises/two-fer"

	for href := range queue {
		webCrawler(href)
	}

}

func webCrawler(href string) {
	fmt.Printf("Crawling URL %v -- \n", href)

	res, err := netClient.Get(href)
	checkerr(err)

	links, err := extractlinks.All(res.Body)
	checkerr(err)

	for _, link := range links {
		//webCrawler(fixURL(link.Href, href))
		absurl := fixURL(link.Href, href)
		go func() {
			queue <- absurl
		}()
	}

	res.Body.Close()
}

func fixURL(href, baseUrl string) string {
	ur, err := url.Parse(href)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}

	fixedRes := base.ResolveReference(ur)

	return fixedRes.String()
}
func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
	}

}
