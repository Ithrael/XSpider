package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/ithrael/WebScan/middleware"
)

var visitURL string

var config *Config
var configFile string

func initConfig() {
	flag.StringVar(&visitURL, "url", "https://www.baidu.com", "URL to visit")
	flag.StringVar(&configFile, "config", "./config.yaml", "config file path")
	flag.Parse()

	config, _ = LoadConfig(configFile)
	if config == nil {
		return
	}
}

func parseResp(resp *colly.Response, url *url.URL) (*middleware.PageDetail, error) {
	html := string(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("Failed to create goquery doc: %v", err)
	}

	title := doc.Find("title").Text()

	h := md5.New()
	h.Write([]byte(html))
	fingerprint := fmt.Sprintf("%x", h.Sum(nil))
	currentTime := time.Now()
	// 格式化为 "2006-01-02 15:04:05" 格式
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return &middleware.PageDetail{
		Url:          url.String(),
		Host:         url.Host,
		Title:        title,
		HTML:         html,
		ResponseCode: resp.StatusCode,
		Fingerprint:  fingerprint,
		Timestamp:    formattedTime,
	}, nil
}

func runSpider(detailsCh chan *middleware.PageDetail) {
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

		if len(config.Restriction.AllowedDomains) != 0 && !IsSubDomain(host, config.Restriction.AllowedDomains) {
			r.Abort()
		}
		if len(config.Restriction.ExcludedDomains) != 0 && IsSubDomain(host, config.Restriction.ExcludedDomains) {
			r.Abort()
		}
		if len(config.Restriction.AllowedPaths) != 0 && !IsRegexMatch(path, config.Restriction.AllowedPaths) {
			r.Abort()
		}
		if len(config.Restriction.ExcludedPaths) != 0 && IsRegexMatch(path, config.Restriction.ExcludedPaths) {
			r.Abort()
		}
		if len(config.Restriction.AllowedQueryKey) != 0 && !IsRegexMatch(queryKey, config.Restriction.AllowedQueryKey) {
			r.Abort()
		}
		if len(config.Restriction.ExcludedQueryKey) != 0 && IsRegexMatch(queryKey, config.Restriction.ExcludedQueryKey) {
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
		pageDetail, err := parseResp(r, r.Request.URL)
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
	initConfig()
	detailsCh := make(chan *middleware.PageDetail)
	// 将数据写入到csv文件(output.csv)中
	go middleware.GetCsvInstance().Process(detailsCh)
	// 将数据写入到mysql中, export MysqlUrl="xxxx"
	if os.Getenv("MYSQL_ENABLE") == "True" {
		go middleware.GetMySQLDBInstance().Process(detailsCh)
	}
	runSpider(detailsCh)
}
