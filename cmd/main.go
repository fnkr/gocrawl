package main

import (
	"flag"
	"github.com/fnkr/gocrawl"
	collyqueue "github.com/gocolly/colly/queue"
	"regexp"
	"os"
)

type regexprsFlag []*regexp.Regexp

func (d *regexprsFlag) String() string {
	// Not actually in use, but needs to be implemented
	return "regexprsFlag.String()"
}

func (d *regexprsFlag) Set(value string) error {
	regex, err := regexp.Compile(value)
	if err != nil {
		return err
	}
	*d = append(*d, regex)
	return nil
}

func main() {
	url := flag.String("url", "", "URL to crawl")
	post := flag.String("post", "", "use POST method for initial request, value will be used as raw post data")
	var allowedURLPatterns regexprsFlag
	flag.Var(&allowedURLPatterns, "only", "only crawl matching URLs, regular expression, e.g. \"^http:\\/www\\.example.com\\/\", can be used multiple times")
	var disallowedURLPatterns regexprsFlag
	flag.Var(&disallowedURLPatterns, "ignore", "prevent specific URLs from being crawled, regular expression, e.g. \"^http:\\/www\\.example.com\\/\", can be used multiple times")
	noParent := flag.Bool("no-parent", false, "do not ascend to the parent directory")
	threads := flag.Int("threads", 1, "number of threads")

	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	crawler := gocrawl.Crawler{
		AllowedURLPatterns:    allowedURLPatterns,
		DisallowedURLPatterns: disallowedURLPatterns,
		NoParent:              *noParent,
		Threads:               *threads,
		InMemoryQueueStorage:  &collyqueue.InMemoryQueueStorage{MaxSize: 100000},
	}
	crawler.Crawl(*url, *post)
}
