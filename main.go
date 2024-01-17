package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"flag"
	"fmt"
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

func isSubdomainOfAllowedDomain(host string, allowedDomains map[string]struct{}) bool {
	for domain := range allowedDomains {
		if strings.HasSuffix(host, domain) {
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

func writeDetailsToCSV(detailsCh chan *PageDetail) error {
	writer, f, err := createCsvWriter()
	if err != nil {
		return err
	}
	defer f.Close()
	defer writer.Flush()

	for detail := range detailsCh {
		writer.Write([]string{
			detail.Url,
			detail.Title,
			fmt.Sprint(detail.ResponseCode),
			detail.Fingerprint,
		})
	}
	return nil
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

func runSpider(detailsCh chan *PageDetail) error {
	c := colly.NewCollector()

	if parallelism != 0 {
		c.Limit(&colly.LimitRule{Parallelism: parallelism, RandomDelay: time.Duration(randomDelayMaxTime) * time.Second})
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		if !isSubdomainOfAllowedDomain(r.URL.Hostname(), allowedDomains) {
			r.Abort()
		} else {
			fmt.Println("Visiting", r.URL)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		pageDetail, err := parseResp(r, r.Request.URL.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		detailsCh <- pageDetail
	})

	if err := c.Visit(visitURL); err != nil {
		return fmt.Errorf("visit url failed: %v", err)
	}
	close(detailsCh)
	return nil
}

var allowedDomains map[string]struct{}
var visitURL string
var out string
var parallelism int
var randomDelayMaxTime int

func init() {
	flag.StringVar(&visitURL, "url", "https://www.baidu.com", "URL to visit")
	flag.StringVar(&out, "out", "out.csv", "out file name")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	randomDelayMaxTime = viper.GetInt("RandomDelayMaxTime")
	parallelism = viper.GetInt("Parallelism")
	domains := viper.GetStringSlice("AllowedDomains")
	allowedDomains = make(map[string]struct{})
	for _, d := range domains {
		allowedDomains[d] = struct{}{}
	}
}

func main() {
	detailsCh := make(chan *PageDetail)
	go writeDetailsToCSV(detailsCh)
	if err := runSpider(detailsCh); err != nil {
		panic(err)
	}
}
