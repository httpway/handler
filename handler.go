package sample

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

// Client has the information needed to work as a httpway plugin.
type Client struct {
	url *url.URL
}

// Close the client.
func (Client) Close(ctx context.Context) error {
	return nil
}

// Default is a sample HTTP handler that track the time of a request.
func (c Client) Default(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Host = c.url.Host
		r.URL.Host = c.url.Host
		r.URL.Scheme = c.url.Scheme

		t := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("request duration: %s\n", time.Since(t).String())
	}
	return http.HandlerFunc(fn)
}

// Panic is used to recover from panics to avoid the server shuting down.
func (Client) Panic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rv := recover(); rv != nil {
				fmt.Fprintf(os.Stderr, "panic: %+v\n", rv)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NotFound is invoked to handle not found errors.
func (Client) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (c *Client) init(config map[string]interface{}) error {
	rawHost, ok := config["host"].(string)
	if !ok {
		return errors.New("casting host to string error")
	}

	var err error
	c.url, err = url.Parse(rawHost)
	if err != nil {
		return errors.Wrapf(err, "parse url '%s' error", rawHost)
	}

	return nil
}

// NewClient return a initialized client.
func NewClient(config map[string]interface{}) (Client, error) {
	var c Client
	if err := c.init(config); err != nil {
		return c, errors.Wrap(err, "initialization error")
	}
	return c, nil
}
