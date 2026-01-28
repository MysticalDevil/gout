package gout

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// H is a shortcut for map[string]any
type H map[string]any

// Context holds the HTTP request and response objects, as well as URL path, method middleware handlers.
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next executes the remaining handlers in the chain.
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Param returns the value of the path parameter specified by the key.
func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

// PostFrom returns the value of the form parameter specified by the key.
func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key)
}

// Query returns the value of the query parameter specified by the key.
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status sets the HTTP response status code.
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader sets the value of the response header.
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String writes a string response with specified status code and format.
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	fmt.Fprintf(c.Writer, format, values...)
}

// JSON writes a JSON response with the specified status code.
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data writes a raw byte slice response with the specified status code.
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML renders a HTML template with the specified status code and data.
func (c *Context) HTML(code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

// ClientIP returns the client's IP address.
func (c *Context) ClientIP() string {
	if ip := c.Req.Header.Get("X-Forwarded-For"); ip != "" {
		if before, _, found := strings.Cut(ip, ","); found {
			return before
		}
		return ip
	}
	if ip := c.Req.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	return c.Req.RemoteAddr
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
func (c *Context) SetCookie(
	name, value string,
	maxAge int,
	path, domain string,
	secure, httpOnly bool,
) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// File writes the specified file into the body stream in an efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Req, filepath)
}

// BindJSON parses the request body as JSON and binds it to a struct.
func BindJSON[T any](c *Context) (T, error) {
	var t T
	if c.Req.Body == nil {
		return t, errors.New("request body is empty")
	}
	decoder := json.NewDecoder(c.Req.Body)
	if err := decoder.Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}

// DataResp defines a standard response struct
type DataResp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Success writes a successful JSON response with a standard structure.
func Success[T any](c *Context, data T) {
	resp := DataResp[T]{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	c.JSON(http.StatusOK, resp)
}
