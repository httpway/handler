package handler

import (
	"context"
	"net/http"
	"net/http/httputil"
)

// Client has the information needed to work as a httpway plugin.
type Client struct{}

// Close the client.
func (Client) Close(ctx context.Context) error {
	return nil
}

// Default act as a reverse proxy.
func (Client) Default(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.ServeHTTP(w, r)
}

// NewClient return a initialized client.
func NewClient() Client {
	return Client{}
}
