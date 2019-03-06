package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

// Client has the information needed to work as a httpway plugin.
type Client struct{}

// Close the client.
func (Client) Close(ctx context.Context) error {
	return nil
}

// Default is a sample HTTP handler that track the time of a request.
func (Client) Default(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		fmt.Println(time.Since(t))
	}
	return http.HandlerFunc(fn)
}

// Panic is used to recover from panics to avoid the server shuting down.
func Panic(next http.Handler) http.Handler {
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
	payload, err := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"status": http.StatusNotFound,
			"uri":    r.URL.String(),
		},
	})
	if err != nil {
		fmt.Println(errors.Wrap(err, "response marshal error").Error())
		return
	}

	if _, err := w.Write(payload); err != nil {
		fmt.Println(errors.Wrap(err, "payload response write error").Error())
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// NewClient return a initialized client.
func NewClient() (Client, error) {
	return Client{}, nil
}
