package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type PageDetail struct {
	Url          string
	Title        string
	HTML         string
	ResponseCode int
	Fingerprint  string
}

var visitURL string
var out string
var config *Config

func init() {
	flag.StringVar(&visitURL, "url", "https://www.baidu.com", "URL to visit")
	flag.StringVar(&out, "out", "out.csv", "out file name")
	flag.Parse()

	config, _ = LoadConfig()
	if config == nil {
		return
	}
}

func parseResp(resp *colly.Response, url string) (*PageDetail, error) {
	html := string(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("Failed to create goquery doc: %v", err)
	}

	title := doc.Find("title").Text()

	h := md5.New()
	h.Write([]byte(html))
	fingerprint := fmt.Sprintf("%x", h.Sum(nil))

	return &PageDetail{
		Url:          url,
		Title:        title,
		HTML:         html,
		ResponseCode: resp.StatusCode,
		Fingerprint:  fingerprint,
	}, nil
}

func runSpider(detailsCh chan *PageDetail) {
	var count int
	c := colly.NewCollector(
		colly.MaxDepth(config.Restriction.MaxDepth),
	)

	if config.Restriction.Parallelism != 0 {
		c.Limit(&colly.LimitRule{Parallelism: config.Restriction.Parallelism, RandomDelay: time.Duration(config.Restriction.RandomDelayMaxTime) * time.Second})
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		if config.Restriction.MaxCount != 0 {
			if count > config.Restriction.MaxCount {
				log.Println("Reached max count!")
				r.Abort()  // 如果超过最大数量，则取消请求
				os.Exit(1) // 并退出程序
			}
		}
		host := r.URL.Host
		path := r.URL.Path
		queryKey := r.URL.RawQuery

		if !IsMatch(host, config.Restriction.AllowedDomains) {
			r.Abort()
		}
		if IsMatch(host, config.Restriction.ExcludedDomains) {
			r.Abort()
		}

		if !IsMatch(path, config.Restriction.AllowedPaths) {
			r.Abort()
		}
		if IsMatch(path, config.Restriction.ExcludedPaths) {
			r.Abort()
		}

		if !IsMatch(queryKey, config.Restriction.AllowedQueryKey) {
			r.Abort()
		}
		if IsMatch(queryKey, config.Restriction.ExcludedQueryKey) {
			r.Abort()
		}

		log.Println("Visiting", r.URL)
		count++
		if len(config.Headers) > 0 {
			for key, value := range config.Headers {
				r.Headers.Set(key, value)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		pageDetail, err := parseResp(r, r.Request.URL.String())
		if err != nil {
			log.Printf("Failed to parse response: %v", err)
			return
		}

		detailsCh <- pageDetail
	})

	if err := c.Visit(visitURL); err != nil {
		log.Fatalf("Visit URL failed: %v", err)
	}
	close(detailsCh)
}

func main() {
	detailsCh := make(chan *PageDetail)
	go WriteDetailsToCSV(detailsCh)
	runSpider(detailsCh)
}
