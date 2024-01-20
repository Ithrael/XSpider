package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type PageDetail struct {
	Url          string
	Title        string
	HTML         string
	ResponseCode int
	Fingerprint  string
}

var allowedDomains map[string]struct{}
var excludeDomains map[string]struct{}
var visitURL string
var out string
var maxCount int
var maxDepth int
var parallelism int
var randomDelayMaxTime int
var headers map[string]string

func init() {
	flag.StringVar(&visitURL, "url", "https://www.baidu.com", "URL to visit")
	flag.StringVar(&out, "out", "out.csv", "out file name")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	randomDelayMaxTime = viper.GetInt("Restriction.RandomDelayMaxTime")
	parallelism = viper.GetInt("Restriction.Parallelism")
	allowdomains := viper.GetStringSlice("Restriction.AllowedDomains")
	excludedomains := viper.GetStringSlice("Restriction.ExcludedDomains")
	maxCount = viper.GetInt("Restriction.MaxCount")
	maxDepth = viper.GetInt("Restriction.MaxDepth")
	headers = viper.GetStringMapString("Headers")
	allowedDomains = make(map[string]struct{})
	for _, d := range allowdomains {
		allowedDomains[d] = struct{}{}
	}

	excludeDomains = make(map[string]struct{})
	for _, d := range excludedomains {
		excludeDomains[d] = struct{}{}
	}
}

func isSubdomainOfAllowedDomain(host string, allowedDomains map[string]struct{}) bool {
	for domain := range allowedDomains {
		if strings.HasSuffix(host, domain) {
			return true
		}
	}
	return false
}

func isExcludedDomain(host string, excludedDomains map[string]struct{}) bool {
	for domain := range excludedDomains {
		if host == domain {
			return true
		}
	}
	return false
}

func createCsvWriter() (*csv.Writer, *os.File, error) {
	csvFile, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(bufio.NewWriter(csvFile))
	return writer, csvFile, nil
}

func writeDetailsToCSV(detailsCh chan *PageDetail) {
	writer, f, err := createCsvWriter()
	if err != nil {
		log.Fatalf("Failed to create CSV writer: %v", err)
	}
	defer f.Close()
	defer writer.Flush()

	for detail := range detailsCh {
		err := writer.Write([]string{
			detail.Url,
			detail.Title,
			fmt.Sprint(detail.ResponseCode),
			detail.Fingerprint,
		})
		if err != nil {
			log.Printf("Failed to write data to CSV: %v", err)
		}
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
		colly.MaxDepth(maxDepth),
	)

	if parallelism != 0 {
		c.Limit(&colly.LimitRule{Parallelism: parallelism, RandomDelay: time.Duration(randomDelayMaxTime) * time.Second})
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		if maxCount != 0 {
			if count > maxCount {
				fmt.Println("Reached max count!")
				r.Abort()  // 如果超过最大数量，则取消请求
				os.Exit(1) // 并退出程序
			}
		}

		if isExcludedDomain(r.URL.Hostname(), excludeDomains) {
			r.Abort()
		}
		if !isSubdomainOfAllowedDomain(r.URL.Hostname(), allowedDomains) {
			r.Abort()
		} else {
			log.Println("Visiting", r.URL)
			count++
			if len(headers) > 0 {
				for key, value := range headers {
					r.Headers.Set(key, value)
				}
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
	go writeDetailsToCSV(detailsCh)
	runSpider(detailsCh)
}
