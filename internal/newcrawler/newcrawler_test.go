package newcrawler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

func TestCrawler_Run(t *testing.T) {
	cfg := dependencies.LoadEnv()

	t.Run("when the starting URL has valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(cfg, f)
		c.Run("https://monzo.com/", 1)

		require.NotEmpty(t, c.visited)

		i := 0
		urls := make([]string, len(c.visited))
		for k := range c.visited {
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
		c.Run("http://dummysite.com/", 1)

		assert.Len(t, c.visited, 1)
		assert.True(t, c.visited["http://dummysite.com/"])
	})

	//t.Run("when crawl depth is limited", func(t *testing.T) {
	//	cfg.MaxCrawlDepth = 2
	//
	//	f := fetcher.NewMockFetcher()
	//	c := NewCrawler(cfg, f)
	//	c.Run("https://monzo.com/", 1)
	//
	//	require.NotEmpty(t, c.visited)
	//
	//	i := 0
	//	urls := make([]string, len(c.visited))
	//	for k := range c.visited {
	//		urls[i] = k
	//		i++
	//	}
	//
	//	assert.ElementsMatch(t, []string{
	//		"https://monzo.com/",
	//		"https://monzo.com/monzo-plus/",
	//		"https://monzo.com/current-account/",
	//	}, urls)
	//})
}
