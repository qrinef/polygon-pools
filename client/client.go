package client

import (
	"fmt"
	"net/http"
	"polygon-pools/reserves"
	"time"
)

type Client struct {
	newReserves *reserves.Reserves
	port        string
}

func NewClient(newReserves *reserves.Reserves) *Client {
	return &Client{
		newReserves: newReserves,
		port:        ":8035",
	}
}

func (c *Client) Start() (err error) {
	fmt.Printf("> Server on %v started!", c.port)

	http.HandleFunc("/reserves", c.handlerReserves)

	s := &http.Server{
		Addr:         c.port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	return s.ListenAndServe()
}

func (c *Client) handlerReserves(rw http.ResponseWriter, _ *http.Request) {
	pools, err := c.newReserves.GetReserves()
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Write(pools)
}
