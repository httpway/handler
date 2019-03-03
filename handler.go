package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client has the information needed to work as a httpway plugin.
type Client struct{}

// Close the client.
func (Client) Close(ctx context.Context) error {
	return nil
}

// Default act as a reverse proxy.
func (c Client) Default(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		fmt.Println(time.Since(t))
	}
	return http.HandlerFunc(fn)
}

// NewClient return a initialized client.
func NewClient() (Client, error) {
	return Client{}, nil
}
