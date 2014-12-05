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
// the request URL before attempting to serve the request
func (f *FileServer) StripPrefix(prefix string)
```

``` go
/// FileServer returns a handler that serves HTTP requests
// with the contents of the directory rooted at root
// calling each option on the server before returning it
func New(root string, options ...func(*FileServer)) http.Handler 
```

package main

import (
	"github.com/qrtz/fs"
	"html/template"
	"net/http"
)

func main() {
	notfound := template.Must(template.New("404").Parse(`<!DOCTYPE html><html><head><title>[404]</title></head><body><h1>[404]</h1></body></html>`))

	http.ListenAndServe(":8080", fs.New(".", func(f *fs.FileServer) {
		f.ErrorHandler(0, func(w http.ResponseWriter, r *http.Request, code int) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			notfound.Execute(w, nil)
		})
	}))
}

```
