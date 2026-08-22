package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL        = "https://www.examtopics.com"
	maxRetries     = 3
	requestsPerSec = 2.0
	httpTimeout    = 25 * time.Second
)

// Question chứa toàn bộ dữ liệu một câu hỏi
type Question struct {
	Title    string   `json:"title"`
	Topic    string   `json:"topic,omitempty"`
	Content  string   `json:"content"`
	Choices  []string `json:"choices"`
	Answer   string   `json:"answer"`
	Time     string   `json:"timestamp,omitempty"`
	URL      string   `json:"url"`
	Comments string   `json:"comments,omitempty"`
}

// Scraper quản lý HTTP client + rate limiter
type Scraper struct {
	client  *http.Client
	workers int
	limiter *time.Ticker
}

func NewScraper(workers int) *Scraper {
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	interval := time.Duration(float64(time.Second) / requestsPerSec)
	return &Scraper{
		client: &http.Client{
			Timeout:   httpTimeout,
			Transport: tr,
		},
		workers: workers,
		limiter: time.NewTicker(interval),
	}
}

// fetch tải URL với retry + exponential backoff
func (s *Scraper) fetch(url string) ([]byte, error) {
	backoff := time.Second
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			jitter := time.Duration(rand.Intn(600)) * time.Millisecond
			time.Sleep(backoff + jitter)
			backoff *= 2
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil
		}
		if resp.StatusCode != http.StatusServiceUnavailable && resp.StatusCode != http.StatusTooManyRequests {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
		}
	}
	return nil, fmt.Errorf("hết lần thử lại cho: %s", url)
}

func (s *Scraper) parseHTML(url string) (*goquery.Document, error) {
	body, err := s.fetch(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}

// GetExams trả về danh sách URL đề thi của một nhà cung cấp
func (s *Scraper) GetExams(provider string) []string {
	doc, err := s.parseHTML(fmt.Sprintf("%s/exams/%s/", baseURL, provider))
	if err != nil {
		return nil
	}
	var exams []string
	doc.Find(".popular-exam-link").Each(func(_ int, sel *goquery.Selection) {
		if href, ok := sel.Attr("href"); ok && href != "" {
			exams = append(exams, baseURL+strings.TrimSpace(href))
		}
	})
	return exams
}

// numPages đọc tổng số trang thảo luận của nhà cung cấp
func (s *Scraper) numPages(provider string) int {
	doc, err := s.parseHTML(fmt.Sprintf("%s/discussions/%s/", baseURL, provider))
	if err != nil {
		return 1
	}
	n := 0
	doc.Find(".discussion-list-page-indicator strong").Each(func(i int, sel *goquery.Selection) {
		if i == 1 {
			n, _ = strconv.Atoi(strings.TrimSpace(sel.Text()))
		}
	})
	if n == 0 {
		return 1
	}
	return n
}

// linksFromPage lấy tất cả href trên 1 trang khớp với grep
func (s *Scraper) linksFromPage(provider string, page int, grep string) []string {
	<-s.limiter.C
	url := fmt.Sprintf("%s/discussions/%s/%d", baseURL, provider, page)
	doc, err := s.parseHTML(url)
	if err != nil {
		return nil
	}
	low := strings.ToLower(grep)
	var links []string
	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")
		if strings.Contains(href, "/discussions/") && strings.Contains(strings.ToLower(href), low) {
			links = append(links, href)
		}
	})
	return links
}

// scrapeQuestion lấy dữ liệu một trang câu hỏi
func (s *Scraper) scrapeQuestion(href string) *Question {
	<-s.limiter.C
	url := baseURL + href
	doc, err := s.parseHTML(url)
	if err != nil {
		return nil
	}

	// Lấy các lựa chọn đáp án
	var choices []string
	doc.Find("li.multi-choice-item").Each(func(_ int, sel *goquery.Selection) {
		if txt := cleanText(sel.Text()); txt != "" {
			choices = append(choices, txt)
		}
	})

	// Lấy đáp án từ JSON embedded (cách mới của ExamTopics)
	answer := extractAnswer(doc)

	return &Question{
		Title:    cleanText(doc.Find("h1").Text()),
		Topic:    strings.ReplaceAll(strings.TrimSpace(doc.Find(".question-discussion-header").Text()), "\t", ""),
		Content:  cleanText(doc.Find(".card-text").Text()),
		Choices:  choices,
		Answer:   answer,
		Time:     cleanText(doc.Find(".discussion-meta-data > i").Text()),
		URL:      url,
		Comments: cleanText(doc.Find(".discussion-container").Text()),
	}
}

// extractAnswer đọc đáp án từ JSON tally hoặc fallback về .correct-answer
func extractAnswer(doc *goquery.Document) string {
	jsonTxt := strings.TrimSpace(doc.Find(".voted-answers-tally script").First().Text())
	if jsonTxt != "" {
		var votes []struct {
			VotedAnswers string `json:"voted_answers"`
			IsMostVoted  bool   `json:"is_most_voted"`
		}
		if json.Unmarshal([]byte(jsonTxt), &votes) == nil {
			for _, v := range votes {
				if v.IsMostVoted {
					return v.VotedAnswers
				}
			}
			if len(votes) > 0 {
				return votes[0].VotedAnswers
			}
		}
	}
	// Fallback: cách cũ
	raw := strings.TrimSpace(doc.Find(".correct-answer").Text())
	if len(raw) > 0 {
		clean := strings.ReplaceAll(strings.ReplaceAll(raw, " ", ""), "\n", "")
		if len(clean) > 0 {
			return string(clean[0])
		}
	}
	return ""
}

// ScrapeAll là hàm chính: scrape tất cả câu hỏi khớp với provider + grep
func (s *Scraper) ScrapeAll(provider, grep string) []Question {
	total := s.numPages(provider)
	fmt.Printf("[*] Tổng số trang thảo luận: %d\n", total)

	// Giai đoạn 1: thu thập links song song
	sem := make(chan struct{}, s.workers)
	var mu sync.Mutex
	var allLinks []string
	var wg sync.WaitGroup
	var donePages int32

	for p := 1; p <= total; p++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(page int) {
			defer func() { <-sem; wg.Done() }()
			links := s.linksFromPage(provider, page, grep)
			mu.Lock()
			allLinks = append(allLinks, links...)
			mu.Unlock()
			cur := atomic.AddInt32(&donePages, 1)
			fmt.Printf("\r    Thu thập links: %d/%d trang  ", cur, total)
		}(p)
	}
	wg.Wait()
	fmt.Println()

	allLinks = dedupe(allLinks)
	allLinks = sortLinks(allLinks)
	fmt.Printf("[*] Tìm thấy %d câu hỏi phù hợp với '%s'.\n", len(allLinks), grep)

	if len(allLinks) == 0 {
		return nil
	}

	// Giai đoạn 2: tải từng câu hỏi song song
	results := make([]*Question, len(allLinks))
	var doneQ int32

	for i, href := range allLinks {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, h string) {
			defer func() { <-sem; wg.Done() }()
			results[idx] = s.scrapeQuestion(h)
			cur := atomic.AddInt32(&doneQ, 1)
			fmt.Printf("\r    Tải câu hỏi: %d/%d  ", cur, len(allLinks))
		}(i, href)
	}
	wg.Wait()
	fmt.Println()

	// Lọc nil
	final := make([]Question, 0, len(results))
	for _, q := range results {
		if q != nil && q.Title != "" {
			final = append(final, *q)
		}
	}
	return final
}

// --- Helpers ---

var wsRe = regexp.MustCompile(`\s+`)

func cleanText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = wsRe.ReplaceAllString(s, " ")
	s = strings.Replace(s, "Suggested Answer", "\nSuggested Answer", 1)
	return strings.TrimSpace(s)
}

func dedupe(links []string) []string {
	seen := make(map[string]struct{}, len(links))
	out := make([]string, 0, len(links))
	for _, l := range links {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			out = append(out, l)
		}
	}
	return out
}

var (
	questionRe = regexp.MustCompile(`question-(\d+)`)
	topicRe    = regexp.MustCompile(`topic-(\d+)`)
)

func extractNum(s string, re *regexp.Regexp) int {
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

func sortLinks(links []string) []string {
	sort.Slice(links, func(i, j int) bool {
		ti, tj := extractNum(links[i], topicRe), extractNum(links[j], topicRe)
		if ti != tj {
			return ti < tj
		}
		return extractNum(links[i], questionRe) < extractNum(links[j], questionRe)
	})
	return links
}
