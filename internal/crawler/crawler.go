package crawler

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

type Crawler struct {
	cfg     *dependencies.Config
	fetcher fetcher.IFetcher
	Visited map[string]bool
	lock    sync.Mutex
}

func NewCrawler(cfg *dependencies.Config, fetcher fetcher.IFetcher) *Crawler {
	return &Crawler{
		cfg:     cfg,
		fetcher: fetcher,
		Visited: make(map[string]bool),
	}
}

// This function simply recurses through parsed links and spins up a goroutine for each new link to visit/crawl.
// It'll spin up as many goroutines as possible to work on each link.
func (c *Crawler) RunUnbounded(url string, depth int) {
	o := c.markAsVisited(url)
	if !o || (c.cfg.MaxCrawlDepth > 0 && depth >= c.cfg.MaxCrawlDepth) {
		return
	}

	log.Printf("visited: %s\n", url)

	urls, err := c.fetcher.Fetch(url)
	if err != nil {
		log.Printf("skipping - unable to crawl %s - %v\n", url, err)
		return
	}

	if len(urls) == 0 {
		return
	}

	c.logAttempts(urls)

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			c.RunUnbounded(u, depth+1)
		}(u)
	}
	wg.Wait()

	return
}

// This function sets a bounded limit on the amount of concurrent web-crawlers that can run at a time.
// It uses a fan-in/fan-out approach by fanning out workers to parse links from concurrent HTTP requests and a main worker to queue up pending links that are waiting to be visited.
func (c *Crawler) RunBounded(url string, depth int) {
	type crawlJob struct {
		url   string
		depth int
	}

	type pendingJob struct {
		urls  []string
		depth int
	}

	var (
		wg   sync.WaitGroup
		work atomic.Int64
	)

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
					continue loop
				}

				log.Printf("visited: %s\n", job.url)

				urls, err := c.fetcher.Fetch(job.url)
				if err != nil {
					log.Printf("skipping - unable to crawl %s - %v\n", job.url, err)
					work.Add(-1)
					continue loop
				}

				if len(urls) == 0 {
					work.Add(-1)
					continue loop
				}

				c.logAttempts(urls)

				pendingUrlsCh <- &pendingJob{urls: urls, depth: job.depth + 1}

				work.Add(-1)
			case <-ctx.Done():
				break loop
			}
		}
	}

	crawl := func(terminator context.CancelFunc, targetUrlCh chan<- *crawlJob, pendingUrlsCh <-chan *pendingJob) {
	loop:
		for {
			select {
			case job := <-pendingUrlsCh:
				for _, u := range job.urls {
					targetUrlCh <- &crawlJob{url: u, depth: job.depth}
				}
				continue loop
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

	pendingUrlsCh := make(chan *pendingJob, 100_000) // Buffered channel to limit the no. of pending unprocessed links at a time.
	defer close(pendingUrlsCh)

	ctx, terminator := context.WithCancel(context.Background())

	for i := 0; i < c.cfg.MaxCrawlConcurrencyLevel; i++ {
		wg.Add(1)
		go worker(ctx, targetUrlCh, pendingUrlsCh)
	}

	// Optional: We could potentially make this function run concurrently as well if we wanted to optimise further.
	go crawl(terminator, targetUrlCh, pendingUrlsCh)

	pendingUrlsCh <- &pendingJob{urls: []string{url}, depth: depth}

	wg.Wait()
}

func (c *Crawler) markAsVisited(url string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.Visited[url]
	if ok {
		return false
	}
	c.Visited[url] = true

	return true
}

func (c *Crawler) logAttempts(urls []string) {
	var visitingUrls string
	if len(urls) > c.cfg.MaxLoggedUrls {
		visitingUrls = fmt.Sprintf("%d links", len(urls))
	} else {
		visitingUrls = strings.Join(urls, ", ")
	}
	log.Printf("\twill try visiting: %s\n", visitingUrls)
}
