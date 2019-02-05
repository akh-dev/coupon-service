package util

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// Since http.ResponseWriter is an interface, we can intercept it and provide our own ResponseWriter to add gzip functionality

// Since gzip is a stream-based protocol, we need to be able to tell when the stream is done so it can flush its buffers
// CloseableResponseWriter is an interface as we might provide different ResponseWriters depending whether the client supports gzip or not
type CloseableResponseWriter interface {
	http.ResponseWriter
	Close()
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

// Since both http.ResponseWriter and gzip.Writer have Write method, we need to do some disambiguation and tell go explicitly which method to invoke
func (w gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// CloseableResponseWriter interface implementation
func (w gzipResponseWriter) Close() {
	w.Writer.Close()
}

// Same as with Write() we need to disambiguate Header() between http.ResponseWriter and gzip.Writer
func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// This is just a wrapper around http.ResponseWriter so we can work with it via CloseableResponseWriter interface
type closeableResponseWriter struct {
	http.ResponseWriter
}

// CloseableResponseWriter interface implementation
func (w closeableResponseWriter) Close() {}

// If the request header "Accept-Encoding" contains "gzip", we will create a gzipResponseWriter and populate its Writer field with a gzip.NewWriter(w)
// Otherwise we will return a standard http.ResponseWriter wrapped into a closeableResponseWriter struct
func GetResponseWriter(w http.ResponseWriter, req *http.Request) CloseableResponseWriter {
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gRW := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzip.NewWriter(w),
		}
		return gRW
	} else {
		return closeableResponseWriter{ResponseWriter: w}
	}
}

//providing our own handler that will be able to serve requests and decide on the fly wich writer to use
type GzipHandler struct{}

//GzipHandler is the only bit that we expose for external use, everything else will be used internally
//(GetResponseWriter is also exposed just in case there is a need for it from outside of the package)
func (h *GzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseWriter := GetResponseWriter(w, r)
	defer responseWriter.Close()

	http.DefaultServeMux.ServeHTTP(responseWriter, r)
}
