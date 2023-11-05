package main

import (
	"flag"
	"log"
	"runtime"
	"time"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
	"webcrawler-go/internal/newcrawler"
)

func main() {
	cfg := dependencies.LoadEnv()
	arg := flag.String("targetUrl", "", "the starting URL that the web-crawler should crawl from.")
	flag.Parse()

	if *arg == "" {
		log.Fatal("web-crawler needs a starting URL")
	}

	// 👋 Enable for benchmarking purposes
	t := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-t:
				go func() {
					log.Printf("No. of goroutines running: %d\n", runtime.NumGoroutine())
				}()
			}
		}
	}()

	start := time.Now()

	f := fetcher.NewFetcher()
	c := newcrawler.NewCrawler(cfg, f)
	//c := crawler.NewCrawler(cfg, f)
	c.Run(*arg, 1)

	end := time.Now()

	log.Printf("✅ web-crawler took %v to complete.\n", end.Sub(start))
}
