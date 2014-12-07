fs
==

Package fs is a simple wrapper for net/http.FileServer. Provides helper functions to configure the behavior of the server

Features
========
* Set Custom Error Handlers
* Set Custom Index pages
* Enable or disable directory listing. Default disabled


``` go
// Set http error handler
// handler with code=0 is used as default handler
func (f *FileServer) ErrorHandler(code int, fn ErrorHandlerFunc)
```

``` go
// Set autoindex to indicate whether to allow directory listing
func (f *FileServer) AutoIndex(autoIndex bool)
```
``` go
// Set index page. Defaults to "index.html"
func (f *FileServer) Index(index ...string)
```
``` go
// Set a handler that remove the prefix from
// the request URL before attempting to serve the request
func (f *FileServer) StripPrefix(prefix string)
```

``` go
/// FileServer returns a handler that serves HTTP requests
// with the contents of the directory rooted at root
// calling each option on the server before returning it
func New(root string, options ...func(*FileServer)) http.Handler 
```

Usage
=====
``` go
package main

import (
	"flag"
	"github.com/qrtz/fs"
	"html/template"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "Server HTTP address")
	root := flag.String("r", ".", "Base directory path for static contents")

	flag.Parse()

	notfound := template.Must(template.New("404").Parse(`<!DOCTYPE html><html><head><title>[404]</title></head><body><h1>[404]</h1></body></html>`))

	http.ListenAndServe(*addr, fs.New(*root, func(f *fs.FileServer) {
		// Set the default error handler
		f.ErrorHandler(0, func(w http.ResponseWriter, r *http.Request, code int) {
			// We can set the content type base on the request
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			// We control the error code to return
			w.WriteHeader(code)

			// We can now respond with our awesome error page.
			notfound.Execute(w, nil)
		})

		// Turn directory listing on. This is off by default
		f.AutoIndex(true)

		// Overwrite the default index page with one or more possible choices
		f.Index("index.html", "default.html", "index.txt")
	}))
}

```
