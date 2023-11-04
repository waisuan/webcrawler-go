package fetcher

import "fmt"

type MockFetcher map[string]*fakeResult

type fakeResult struct {
	urls []string
}

func NewMockFetcher() MockFetcher {
	return MockFetcher{
		"http://golang.org/": &fakeResult{
			[]string{
				"http://golang.org/pkg/",
				"http://golang.org/cmd/",
			},
		},
		"http://golang.org/pkg/": &fakeResult{
			[]string{
				"http://golang.org/",
				"http://golang.org/cmd/",
				"http://golang.org/pkg/fmt/",
				"http://golang.org/pkg/os/",
			},
		},
		"http://golang.org/pkg/fmt/": &fakeResult{
			[]string{
				"http://golang.org/",
				"http://golang.org/pkg/",
			},
		},
		"http://golang.org/pkg/os/": &fakeResult{
			[]string{
				"http://golang.org/",
				"http://golang.org/pkg/",
			},
		},
	}
}

func (f MockFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		return res.urls, nil
	}

	return nil, fmt.Errorf("not found: %s", url)
}
