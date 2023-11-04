package fetcher

type IFetcher interface {
	Fetch(url string) ([]string, error)
}

type Fetcher struct{}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Fetch(url string) ([]string, error) {
	return nil, nil
}
