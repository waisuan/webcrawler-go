# Overview

This is a simple web-crawler built in Go. It takes in one argument (the starting URL) and crawls the site for all links that belong to the same domain. It'll print the links that it visits along the way as well as the links that it'll visit next. The crawler will skip over links that it has visited previously. The crawler also uses a best-effort approach when parsing links and will simply skip over any that can't be reached (e.g. due to HTTP timeouts).

# Pre-requisites

`>= go 1.20`

# Usage

Run app in default mode: `make run targetUrl="<URL>"`

Run app in dev mode: `make dev targetUrl="<URL>"`

Run unit tests:
`make test`

## Environment variables

Add them to their respective `.env` files in order to configure the crawler's behaviour. Refer to `config.go` to view their default values.

`MAX_CRAWL_CONCURRENCY_LEVEL`

Limit the no. of running goroutines. This is useful for also limiting the no. of concurrent HTTP requests made at a time. By default, this value is unbounded.

`MAX_CRAWL_DEPTH`

Limit the depth of pages/links the crawler should process. This is useful for indirectly controlling how long the crawler should run for. By default, this value is unbounded.

`MAX_LOGGED_URLS`

Limit the amount of pending links printed to the console. E.g. "will try visiting: URL1, URL2, ..." -> "will try visiting: 500 links"

# Future State

The following are some of the action items that the developer would like to visit/address if/when time permits that'd help strengthen the quality, resiliency, and observability of the crawler.

- Add external storage (cache/DB) to host all visited links. Can also help with analysing/querying links that were visited on a certain datetime.
- Configure HTTP timeout when fetching HTML pages to avoid waiting too long for a page to respond.
- Configure operational timeout when running the crawler so that it doesn't end up running for an indefinite amount of time.
- Acknowledge site security/privacy settings and explicitly skip over links that should not be visited. E.g. robots.txt
- Better output reporting mechanism -- for querying/analytics purposes. E.g. report failed/skipped/successful links to a persistent storage device.
- Add exponential retries to crawler to handle intermittent runtime errors (e.g. HTTP timeouts).
- Add sleep time in between crawls to be less disruptive towards and reduce load pressure on target sites.
- Benchmark crawler to identify concurrency limits.
- Add linter to enforce code quality.
- Better error reporting mechanism -- for monitoring purpose. E.g. send runtime errors to DataDog where devs can easily build custom alarms around.
- Add custom telemetry around crawler behaviour -- helps to identify unhealthy system anomalies. E.g. send custom metrics to DataDog where devs can easily build custom dashboards around.