package fs

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	errNoContentWritten = errors.New("content not written")
	errForbidden        = errors.New("forbidden")
)

type ResponseWrapper struct {
	http.ResponseWriter
	statusCode int
	buf        []byte
}

func (w *ResponseWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWrapper) Write(b []byte) (int, error) {
	if w.statusCode < http.StatusBadRequest {
		return w.ResponseWriter.Write(b)
	}
	w.buf = b
	return 0, errNoContentWritten
}

func (w *ResponseWrapper) WriteHeader(code int) {
	w.statusCode = code
	if code < http.StatusBadRequest {
		w.ResponseWriter.WriteHeader(code)
	}
}

type Dir struct {
	autoIndex bool
	index     []string
	err       error
	http.FileSystem
}

func (d *Dir) Open(name string) (http.File, error) {
	f, err := d.FileSystem.Open(name)

	if err != nil {
		d.err = err
		return nil, err
	}

	info, err := f.Stat()

	if err != nil {
		d.err = err
		return nil, err
	}

	if info.IsDir() {
		for _, index := range d.index {
			if ff, err := d.FileSystem.Open(filepath.Join(name, index)); err == nil {
				f.Close()
				return ff, nil
			}
		}

		if !d.autoIndex {
			f.Close()
			d.err = errForbidden
			return nil, d.err
		}
	}

	return f, nil
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, int)

type FileServer struct {
	errorHandler map[int]ErrorHandlerFunc
	autoIndex    bool     // Indicates whether to allow directory listing
	index        []string // list of index pages
	prefix       string
	root         string
	handlerFunc  func(f *FileServer, w *ResponseWrapper, r *http.Request)
}

// Set http error handler
// handler with code=0 is used as default handler
func (f *FileServer) ErrorHandler(code int, fn ErrorHandlerFunc) {
	f.errorHandler[code] = fn
}

// Set autoindex to indicate whether to allow directory listing
func (f *FileServer) AutoIndex(autoIndex bool) {
	f.autoIndex = autoIndex
}

// Set index page. Defaults to "index.html"
func (f *FileServer) Index(index ...string) {
	f.index = index
}

// Set a handler that remove the prefix from
// the request URL before attempting to server request
func (f *FileServer) StripPrefix(prefix string) {
	f.prefix = prefix
	handlerFunc := f.handlerFunc
	f.handlerFunc = func(s *FileServer, w *ResponseWrapper, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, s.prefix)

		if len(p) == len(r.URL.Path) {
			http.NotFound(w, r)
			s.handleError(w, r)
			return
		}

		r.URL.Path = p
		handlerFunc(s, w, r)
	}
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := filepath.Base(r.URL.Path)
	// Redirect to ../[index] to /. Similar to the behavior of http.FileServer
	for _, index := range f.index {
		if index == name {
			path := filepath.Dir(r.URL.Path)
			if q := r.URL.RawQuery; q != "" {
				path += "?" + q
			}
			w.Header().Set("Location", path)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	f.handlerFunc(f, &ResponseWrapper{ResponseWriter: w}, r)
}

func (f *FileServer) serveHTTP(w *ResponseWrapper, r *http.Request) {
	fileSystem := &Dir{f.autoIndex, f.index, nil, http.Dir(f.root)}

	http.FileServer(fileSystem).ServeHTTP(w, r)

	if w.statusCode >= http.StatusBadRequest {
		if fileSystem.err == errForbidden {
			w.statusCode = http.StatusForbidden
		}

		f.handleError(w, r)
	}
}

func (f *FileServer) handleError(w *ResponseWrapper, r *http.Request) {
	fn, ok := f.errorHandler[w.statusCode]
	if !ok {
		fn = f.errorHandler[0]
	}
	if fn != nil {
		fn(w.ResponseWriter, r, w.statusCode)
	} else {
		http.Error(w.ResponseWriter, string(w.buf), w.statusCode)
	}
}

// FileServer returns a handler that serves HTTP requests
// with the contents of the directory rooted at root
// calling each option on the server before returning it
func New(root string, options ...func(*FileServer)) http.Handler {
	f := &FileServer{
		errorHandler: make(map[int]ErrorHandlerFunc),
		index:        []string{"index.html"},
		root:         root,
		handlerFunc: func(f *FileServer, w *ResponseWrapper, r *http.Request) {
			f.serveHTTP(w, r)
		},
	}

	for _, option := range options {
		option(f)
	}

	return f
}
