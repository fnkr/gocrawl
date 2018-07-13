package gocrawl

import (
	"regexp"
	"github.com/gocolly/colly"
	collyqueue "github.com/gocolly/colly/queue"
	"fmt"
	"net/url"
)

type Crawler struct {
	AllowedURLPatterns    []*regexp.Regexp
	DisallowedURLPatterns []*regexp.Regexp
	NoParent              bool
	Threads               int
	InMemoryQueueStorage  *collyqueue.InMemoryQueueStorage
}

func (crawler *Crawler) Crawl(startURL, post string) {
	collector := colly.NewCollector(
		colly.URLFilters(crawler.AllowedURLPatterns...),
		colly.DisallowedURLFilters(crawler.DisallowedURLPatterns...),
		colly.UserAgent("gocrawl - https://github.com/fnkr/gocrawl"),
	)
	queue, _ := collyqueue.New(crawler.Threads, crawler.InMemoryQueueStorage)
	noParent := crawler.NoParent
	var resolvedStartURL *url.URL

	collector.OnResponse(func(resp *colly.Response) {
		if resolvedStartURL == nil {
			resolvedStartURL = resp.Request.URL
		}
		fmt.Println(resp.StatusCode, resp.Request.Method, resp.Request.URL)
	})

	collector.OnHTML("a[href]", func(link *colly.HTMLElement) {
		resolvedHref := link.Request.AbsoluteURL(link.Attr("href"))

		follow := true

		if noParent {
			if parsedHref, err := url.Parse(resolvedHref); err != nil || isParent(resolvedStartURL, parsedHref) {
				follow = false
			}
		}

		if follow {
			queue.AddURL(resolvedHref)
		}
	})

	if len(post) != 0 {
		collector.PostRaw(startURL, []byte(post))
	} else {
		collector.Visit(startURL)
	}

	queue.Run(collector)
}
