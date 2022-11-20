package api

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}
