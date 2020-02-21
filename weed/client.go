package weed

import "net/http"

var client = http.DefaultClient

func InitHttpClient(c *http.Client) {
	client = c
}

func HttpClient() *http.Client {
	return client
}
