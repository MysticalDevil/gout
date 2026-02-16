package gout

import "testing"

func TestInsertReusesStaticChildren(t *testing.T) {
	root := &node{}

	root.insert("/hello/a", parsePattern("/hello/a"), 0)
	root.insert("/hello/b", parsePattern("/hello/b"), 0)

	if len(root.children) != 1 {
		t.Fatalf("expected one shared static child at root, got %d", len(root.children))
	}
	if root.children[0].part != "hello" {
		t.Fatalf("expected shared child part to be 'hello', got %q", root.children[0].part)
	}
}
