package fetcher

import "fmt"

type MockFetcher map[string]*fakeResult

type fakeResult struct {
	urls []string
}

func NewMockFetcher() MockFetcher {
	return MockFetcher{
		"https://monzo.com/": &fakeResult{
			[]string{
				"https://monzo.com/current-account/",
				"https://monzo.com/monzo-plus/",
			},
		},
		"https://monzo.com/current-account/": &fakeResult{
			[]string{
				"https://monzo.com/",
				"https://monzo.com/help/",
				"https://monzo.com/current-account/joint-account/",
				"https://monzo.com/switch/",
			},
		},
		"https://monzo.com/current-account/joint-account/": &fakeResult{
			[]string{
				"https://monzo.com/",
				"https://monzo.com/current-account/",
			},
		},
		"https://monzo.com/switch/": &fakeResult{
			[]string{
				"https://monzo.com/",
				"https://monzo.com/current-account/",
			},
		},
		"https://monzo.com/monzo-plus/": &fakeResult{
			[]string{},
		},
	}
}

func (f MockFetcher) Fetch(targetUrl string) ([]string, error) {
	if res, ok := f[targetUrl]; ok {
		return res.urls, nil
	}

	return nil, fmt.Errorf("not found: %s", targetUrl)
}
