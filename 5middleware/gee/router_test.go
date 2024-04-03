package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := NewRouter()

	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/b/c", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(ParsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(ParsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(ParsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.GetRoute("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}