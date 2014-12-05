fs
==

Package fs is a simple wrapper for net/http.FileServer. Provides helper functions to configure the behavior of the server

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
// the request URL before attempting to server request
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
	"github.com/qrtz/fileserver"
	"html/template"
	"net/http"
)

// Server files in the current directory
func main() {
	notfound := template.Must(template.New("404").Parse(`<!DOCTYPE html><html><body><h1>[404]</h1></body></html>`))

	http.ListenAndServe(":8080", fileserver.New(".", func(fs *fileserver.Server) {
		fs.ErrorHandler(0, func(w http.ResponseWriter, r *http.Request, code int) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			notfound.Execute(w, nil)
		})
	}))
}

```
