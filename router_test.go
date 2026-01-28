package gout

import (
	"fmt"
	"slices"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/asserts/*filepath", nil)
	return r
}

func TestParsePath(t *testing.T) {
	ok := slices.Equal(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && slices.Equal(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && slices.Equal(parsePattern("/p/*name"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/gout")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "gout" {
		t.Fatal("name should be equal to 'gout'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}
