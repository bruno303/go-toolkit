package http

import (
	"errors"
	"net/http"
)

type Chain struct {
	start Middleware
}

func NewChain(mds ...Middleware) (*Chain, error) {
	chainLen := len(mds)
	if chainLen == 0 {
		return nil, errors.New("chain must not be empty")
	}

	for i := 0; i < chainLen-1; i++ {
		mds[i].SetNext(mds[i+1])
	}

	return &Chain{start: mds[0]}, nil
}

func (c *Chain) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if c.start != nil {
		c.start.ServeHTTP(rw, r)
	}
}
