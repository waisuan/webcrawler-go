package newcrawler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

var (
	wg   sync.WaitGroup
	work atomic.Int64
)

type Crawler struct {
	cfg     *dependencies.Config
	fetcher fetcher.IFetcher
	visited map[string]bool
	lock    sync.Mutex
}

func NewCrawler(cfg *dependencies.Config, fetcher fetcher.IFetcher) *Crawler {
	return &Crawler{
		cfg:     cfg,
		fetcher: fetcher,
		visited: make(map[string]bool),
	}
}

func (c *Crawler) Run(url string, depth int) {
	type crawlJob struct {
		url   string
		depth int
	}

	type pendingJob struct {
		urls  []string
		depth int
	}

	worker := func(ctx context.Context, targetUrlCh <-chan *crawlJob, pendingUrlsCh chan<- *pendingJob) {
		defer wg.Done()

	loop:
		for {
			select {
			case job := <-targetUrlCh:
				work.Add(1)

				o := c.markAsVisited(job.url)
				if !o || (c.cfg.MaxCrawlDepth > 0 && job.depth >= c.cfg.MaxCrawlDepth) {
					work.Add(-1)
					continue
				}

				log.Printf("visited: %s\n", job.url)

				urls, err := c.fetcher.Fetch(job.url)
				if err != nil {
					log.Printf("skipping - unable to crawl %s - %v\n", job.url, err)
					work.Add(-1)
					continue
				}

				if len(urls) == 0 {
					work.Add(-1)
					continue
				}

				var visitingUrls string
				if len(urls) > c.cfg.MaxLoggedUrls {
					visitingUrls = fmt.Sprintf("%d links", len(urls))
				} else {
					visitingUrls = strings.Join(urls, ", ")
				}
				log.Printf("\twill try visiting: %s\n", visitingUrls)

				pendingUrlsCh <- &pendingJob{urls: urls, depth: job.depth + 1}

				work.Add(-1)
			case <-ctx.Done():
				break loop
			}
		}
	}

	crawl := func(terminator context.CancelFunc, targetUrlCh chan<- *crawlJob, pendingUrlsCh chan *pendingJob) {
	loop:
		for {
			select {
			case job := <-pendingUrlsCh:
				for _, u := range job.urls {
					targetUrlCh <- &crawlJob{url: u, depth: job.depth}
				}
				continue
			case <-time.After(3 * time.Second): // helps to terminate all workers when there's nothing left to process.
				log.Printf("searching for more links to process...(%d)\n", work.Load())
				if work.Load() <= 0 {
					break loop
				}
			}
		}

		terminator()
	}

	targetUrlCh := make(chan *crawlJob)
	defer close(targetUrlCh)

	pendingUrlsCh := make(chan *pendingJob, 10_000) // Buffered channel to limit the no. of pending unprocessed links at a time.
	defer close(pendingUrlsCh)

	ctx, terminator := context.WithCancel(context.Background())

	for i := 0; i < c.cfg.MaxCrawlConcurrencyLevel; i++ {
		wg.Add(1)
		go worker(ctx, targetUrlCh, pendingUrlsCh)
	}

	go crawl(terminator, targetUrlCh, pendingUrlsCh)

	pendingUrlsCh <- &pendingJob{urls: []string{url}, depth: depth}

	wg.Wait()
}

func (c *Crawler) markAsVisited(url string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.visited[url]
	if ok {
		return false
	}
	c.visited[url] = true

	return true
}
