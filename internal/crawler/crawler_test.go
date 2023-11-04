package crawler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"webcrawler-go/internal/fetcher"
)

func TestCrawler_Run(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		f := fetcher.NewMockFetcher()
		c := NewCrawler(f)
		c.Run("http://golang.org/")

		require.NotEmpty(t, c.visited)

		i := 0
		urls := make([]string, len(c.visited))
		for k := range c.visited {
			urls[i] = k
			i++
		}

		assert.ElementsMatch(t, []string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		}, urls)
	})
}
