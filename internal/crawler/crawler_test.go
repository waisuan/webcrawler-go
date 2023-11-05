package crawler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"testing"
	"time"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

func TestCrawler_RunUnbounded(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	cfg := dependencies.LoadEnv()

	t.Run("when the starting URL has valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunUnbounded("https://monzo.com/", 1)

		require.NotEmpty(t, c.Visited)

		i := 0
		urls := make([]string, len(c.Visited))
		for k := range c.Visited {
			urls[i] = k
			i++
		}

		assert.ElementsMatch(t, []string{
			"https://monzo.com/",
			"https://monzo.com/current-account/",
			"https://monzo.com/current-account/joint-account/",
			"https://monzo.com/switch/",
			"https://monzo.com/monzo-plus/",
			"https://monzo.com/help/",
		}, urls)
	})

	t.Run("when the starting URL has no valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunUnbounded("http://dummysite.com/", 1)

		assert.Len(t, c.Visited, 1)
		assert.True(t, c.Visited["http://dummysite.com/"])
	})

	t.Run("when crawl depth is limited", func(t *testing.T) {
		cfg.MaxCrawlDepth = 2

		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunUnbounded("https://monzo.com/", 1)

		require.NotEmpty(t, c.Visited)

		i := 0
		urls := make([]string, len(c.Visited))
		for k := range c.Visited {
			urls[i] = k
			i++
		}

		assert.ElementsMatch(t, []string{
			"https://monzo.com/",
			"https://monzo.com/monzo-plus/",
			"https://monzo.com/current-account/",
		}, urls)
	})
}

func TestCrawler_RunBounded(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	cfg := dependencies.LoadEnv()

	t.Run("when the starting URL has valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunBounded("https://monzo.com/", 1)

		require.NotEmpty(t, c.Visited)

		i := 0
		urls := make([]string, len(c.Visited))
		for k := range c.Visited {
			urls[i] = k
			i++
		}

		assert.ElementsMatch(t, []string{
			"https://monzo.com/",
			"https://monzo.com/current-account/",
			"https://monzo.com/current-account/joint-account/",
			"https://monzo.com/switch/",
			"https://monzo.com/monzo-plus/",
			"https://monzo.com/help/",
		}, urls)
	})

	t.Run("when the starting URL has no valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunBounded("http://dummysite.com/", 1)

		assert.Len(t, c.Visited, 1)
		assert.True(t, c.Visited["http://dummysite.com/"])
	})

	t.Run("when crawl depth is limited", func(t *testing.T) {
		cfg.MaxCrawlDepth = 2

		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.RunBounded("https://monzo.com/", 1)

		require.NotEmpty(t, c.Visited)

		i := 0
		urls := make([]string, len(c.Visited))
		for k := range c.Visited {
			urls[i] = k
			i++
		}

		assert.ElementsMatch(t, []string{
			"https://monzo.com/",
			"https://monzo.com/monzo-plus/",
			"https://monzo.com/current-account/",
		}, urls)
	})

	t.Run("respects the max concurrency limit", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		go c.RunBounded("https://monzo.com/", 1)

		time.Sleep(1 * time.Second)

		require.Equal(t, cfg.MaxCrawlConcurrencyLevel+5, runtime.NumGoroutine()) // guesstimate: +5 is the sweet spot reserved for non-worker goroutines operation in this function call.
	})
}
