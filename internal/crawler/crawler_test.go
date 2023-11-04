package crawler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"webcrawler-go/internal/fetcher"
)

func TestCrawler_Run(t *testing.T) {
	t.Run("when the starting URL has valid links", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(f)
		c.Run("https://monzo.com/")

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
		c := NewCrawler(f)
		c.Run("http://dummysite.com/")

		assert.Len(t, c.visited, 1)
		assert.True(t, c.visited["http://dummysite.com/"])
	})
}
