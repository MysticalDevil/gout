package gout

import "net/http"

// CORS returns a middleware handler that enables Cross-Origin Resources Sharing.
func CORS() HandlerFunc {
	return func(c *Context) {
		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.SetHeader(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)
		c.SetHeader("Access-Control-Expose-Headers", "Content-Length")

		if c.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
