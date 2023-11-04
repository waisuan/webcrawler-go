package fetcher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetcher_Fetch(t *testing.T) {
	t.Run("when the HTML page has urls", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `<!doctype html>
<html>
  <head>
    <title>This is the title of the webpage!</title>
  </head>
  <body>
    <p>This is an example paragraph. Anything in the <strong>body</strong> tag will appear on the page, just like this <strong>p</strong> tag and its contents.</p>
	<p><a href="/monzo-plus">Visit!</a></p>
	<p><a href="https://monzo.com/about/">Visit!</a></p>
  </body>
</html>`)
		}))
		defer testServer.Close()

		f := NewFetcher()
		urls, err := f.Fetch("https://monzo.com")
		fmt.Println(urls)
		assert.Nil(t, err)
	})
}
