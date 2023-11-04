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
	<link rel="icon" href="/favicon.png" sizes="any"/>
	<link rel="stylesheet" href="/_next/static/css/90b48bfe6829255e.css" data-n-g=""/>
  </head>
  <body>
    <p>This is an example paragraph. Anything in the <strong>body</strong> tag will appear on the page, just like this <strong>p</strong> tag and its contents.</p>
	<p><a href="/settings">Settings</a></p>
	<p><a href="http://%s/about/">About us</a></p>
	<p><a href="https://monzo.com/about/">Monzo - About Us</a></p>
	<p><a href="https://google.com">Looking for something?</a></p>
  </body>
</html>`, r.Host)
		}))
		defer testServer.Close()

		f := NewFetcher()
		urls, err := f.Fetch(testServer.URL)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{
			fmt.Sprintf("%s/about/", testServer.URL),
			fmt.Sprintf("%s/settings", testServer.URL),
		}, urls)
	})

	t.Run("when the HTML page has no urls", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `<!doctype html>
<html>
  <head>
    <title>This is the title of the webpage!</title>
	<link rel="icon" href="/favicon.png" sizes="any"/>
	<link rel="stylesheet" href="/_next/static/css/90b48bfe6829255e.css" data-n-g=""/>
  </head>
  <body>
    <p>This is an example paragraph. Anything in the <strong>body</strong> tag will appear on the page, just like this <strong>p</strong> tag and its contents.</p>
  </body>
</html>`)
		}))
		defer testServer.Close()

		f := NewFetcher()
		urls, err := f.Fetch(testServer.URL)
		assert.Nil(t, err)
		assert.Empty(t, urls)
	})

	t.Run("when unable to access HTML page", func(t *testing.T) {
		f := NewFetcher()
		urls, err := f.Fetch("https://localhost.org/")
		assert.ErrorContains(t, err, "no such host")
		assert.Empty(t, urls)
	})

	t.Run("when target URL is invalid", func(t *testing.T) {
		f := NewFetcher()
		urls, err := f.Fetch("MALFOMRED_URL.")
		assert.ErrorContains(t, err, "unsupported protocol scheme")
		assert.Empty(t, urls)
	})
}
