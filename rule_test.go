package router

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/novakit/nova"
)

func TestRouteHeaderRule_Match(t *testing.T) {
	r := HeaderRule{Name: "X-Test-Header", Value: []string{"AAA", "BBB"}}
	req := &http.Request{Header: http.Header{}}
	c := &nova.Context{Req: req, Values: map[string]interface{}{}}
	c.Req.Header.Set("X-Test-Header", "AAA")
	if !r.Match(c) {
		t.Error("header rule failed 1")
	}
	c.Req.Header.Set("X-Test-Header", "BBB")
	if !r.Match(c) {
		t.Error("header rule failed 2")
	}
	c.Req.Header.Set("X-Test-Header", "CCC")
	if r.Match(c) {
		t.Error("header rule failed 3")
	}
}

func TestRouteHostRule_Match(t *testing.T) {
	r := HostRule{Host: []string{"host1.landzero.net", "host2.landzero.net"}}
	req := &http.Request{Host: "host1.landzero.net"}
	c := &nova.Context{Req: req, Values: map[string]interface{}{}}
	if !r.Match(c) {
		t.Error("host rule failed 1")
	}
	c.Req.Host = "host2.landzero.net"
	if !r.Match(c) {
		t.Error("host rule failed 2")
	}
	c.Req.Host = "islandzero.net"
	if r.Match(c) {
		t.Error("host rule failed 3")
	}
}

func TestRouteMethodRule_Match(t *testing.T) {
	r := MethodRule{Method: []string{http.MethodGet, http.MethodPost}}
	req := &http.Request{}
	c := &nova.Context{Req: req, Values: map[string]interface{}{}}
	c.Req.Method = http.MethodGet
	if !r.Match(c) {
		t.Error("method rule failed 1")
	}
	c.Req.Method = http.MethodPost
	if !r.Match(c) {
		t.Error("method rule failed 2")
	}
	c.Req.Method = http.MethodPut
	if r.Match(c) {
		t.Error("method rule failed 3")
	}
}

func TestRoutePathRule_Match(t *testing.T) {
	r := PathRule{Path: "/hello/world"}
	req := &http.Request{}
	c := &nova.Context{Req: req, Values: map[string]interface{}{}}
	c.Req.URL, _ = url.Parse("/hello/world")
	if !r.Match(c) {
		t.Error("path rule failed 1")
	}
	c.Req.URL, _ = url.Parse("/hello/world1")
	if r.Match(c) {
		t.Error("path rule failed 2")
	}
	c.Req.URL, _ = url.Parse("/hello/world/1")
	if r.Match(c) {
		t.Error("path rule failed 3")
	}
	r = PathRule{Path: "hello/world/:name"}
	if !r.Match(c) {
		t.Error("path rule failed 4")
	} else {
		if PathParams(c).Get("name") != "1" {
			t.Error("path rule failed 5")
		}
	}
	c.Req.URL, _ = url.Parse("/hello/world/1/bug")
	if r.Match(c) {
		t.Error("path rule failed 6")
	}
	r = PathRule{Path: "hello/world/*any"}
	c.Req.URL, _ = url.Parse("/hello")
	if r.Match(c) {
		t.Error("path rule failed 7")
	}
	c.Req.URL, _ = url.Parse("/hello/world/some/thing/wired")
	if !r.Match(c) {
		t.Error("path rule failed 8")
	} else {
		if PathParams(c).Get("any") != "some/thing/wired" {
			t.Error("path rule failed 9")
		}
	}
	r = PathRule{Path: "hello/:name/*any"}
	c.Req.URL, _ = url.Parse("/hello/world/something/wired")
	if !r.Match(c) {
		t.Error("path rule failed 10")
	} else {
		if PathParams(c).Get("any") != "something/wired" || PathParams(c).Get("name") != "world" {
			t.Error("path rule failed 11")
		}
	}
}

func TestRouteRules_Add(t *testing.T) {
	// nil should work
	var r Rules
	r = r.Add(MethodRule{Method: []string{http.MethodGet, http.MethodPost}})
	req := &http.Request{}
	c := &nova.Context{Req: req, Values: map[string]interface{}{}}
	c.Req.Method = http.MethodGet
	if !r.Match(c) {
		t.Error("add failed 1")
	}
	c.Req.Method = http.MethodPost
	if !r.Match(c) {
		t.Error("add failed 2")
	}
	c.Req.Method = http.MethodPut
	if r.Match(c) {
		t.Error("add failed 3")
	}
	// multiple
	r = r.Add(HostRule{Host: []string{"host1.landzero.net", "host2.landzero.net"}})
	c.Req.Method = http.MethodGet
	c.Req.Host = "host2.landzero.net"
	if !r.Match(c) {
		t.Error("add failed 4")
	}
	c.Req.Method = http.MethodPut
	if r.Match(c) {
		t.Error("add failed 5")
	}
	c.Req.Method = http.MethodGet
	c.Req.Host = "host3.landzero.net"
	if r.Match(c) {
		t.Error("add failed 6")
	}
}

func TestRouteRules_Match(t *testing.T) {
	// already tested
}
