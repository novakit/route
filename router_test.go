package router

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/novakit/nova"
	"github.com/novakit/testkit"
)

func TestRoute_All(t *testing.T) {
	var called bool
	var name string
	n := nova.New()
	Route(n).Put("/hello/:name").Host("landzero.net").Use(func(c *nova.Context) error {
		called = true
		name = PathParams(c).Get("name")
		c.Res.Write([]byte(name))
		return nil
	})
	r := &http.Request{
		Host:   "landzero.net",
		Method: http.MethodPut,
		Header: http.Header{},
	}
	r.URL, _ = url.Parse("/hello/world")
	p := testkit.NewDummyResponse()
	n.ServeHTTP(p, r)
	if !called || name != "world" || p.String() != "world" {
		t.Error("failed")
	}
}
