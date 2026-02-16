package gout

import "testing"

func TestOPTIONSRegistersUnderOptionsMethod(t *testing.T) {
	engine := New()
	engine.OPTIONS("/preflight", func(c *Context) {})

	n, _ := engine.router.getRoute("OPTIONS", "/preflight")
	if n == nil {
		t.Fatal("expected OPTIONS route to be registered")
	}

	putNode, _ := engine.router.getRoute("PUT", "/preflight")
	if putNode != nil {
		t.Fatal("OPTIONS route should not be registered as PUT")
	}
}
