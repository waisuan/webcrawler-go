package main

import (
	"flag"
	"log"
	"time"
	"webcrawler-go/internal/crawler"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

func main() {
	cfg := dependencies.LoadEnv()
	arg := flag.String("targetUrl", "", "the starting URL that the web-crawler should crawl from.")
	flag.Parse()

	if *arg == "" {
		log.Fatal("web-crawler needs a starting URL")
	}

	// 👋 Enable for benchmarking purposes
	//t := time.Tick(time.Second)
	//go func() {
	//	for {
	//		select {
	//		case <-t:
	//			log.Printf("No. of goroutines running: %d\n", runtime.NumGoroutine())
	//		}
	//	}
	//}()

	start := time.Now()

	f := fetcher.NewFetcher()
	c := crawler.NewCrawler(cfg, f)

	if cfg.MaxCrawlConcurrencyLevel > 0 {
		log.Println("Running in BOUNDED mode...")
		c.RunBounded(*arg, 1)
	} else {
		log.Println("Running in UNBOUNDED mode...")
		c.RunUnbounded(*arg, 1)
	}

	end := time.Now()

	log.Printf("✅ web-crawler visited %d links and took %v to complete.\n", len(c.Visited), end.Sub(start))
}
