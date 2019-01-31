package xroute

// The original work was derived from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type ResponseWriterWithStatus interface {
	http.ResponseWriter
	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int
	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int
	HasStatus(status ...int) bool
}

func ResponseWriter(w http.ResponseWriter, status ...int) ResponseWriterWithStatus {
	if ws, ok := w.(ResponseWriterWithStatus); ok {
		return ws
	}
	if len(status) == 0 {
		status = []int{0}
	}
	return &BasicWriter{ResponseWriter:w}
}

type DefaultResponseWriterWithStatus struct {
	http.ResponseWriter
	status int
}

func (w *DefaultResponseWriterWithStatus) Status() int {
	return w.status
}

func (w *DefaultResponseWriterWithStatus) IsHeaderSent() bool {
	return w.status != 0
}

func (rw *DefaultResponseWriterWithStatus) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// WrapResponseWriter is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type WrapResponseWriter interface {
	ResponseWriterWithStatus
	// Tee causes the response body to be written to the given io.Writer in
	// addition to proxying the writes through. Only one io.Writer can be
	// tee'd to at once: setting a second one will overwrite the first.
	// Writes will be sent to the proxy before being written to this
	// io.Writer. It is illegal for the tee'd writer to be modified
	// concurrently with writes.
	Tee(io.Writer)
	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// BasicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type BasicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	bytes       int
	tee         io.Writer
}

func (b *BasicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.status = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}
func (b *BasicWriter) Write(buf []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	n, err := b.ResponseWriter.Write(buf)
	if b.tee != nil {
		_, err2 := b.tee.Write(buf[:n])
		// Prefer errors generated by the proxied writer.
		if err == nil {
			err = err2
		}
	}
	b.bytes += n
	return n, err
}
func (b *BasicWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}
func (b *BasicWriter) Status() int {
	return b.status
}
func (b *BasicWriter) HasStatus(status ...int) bool {
	for _, s := range status {
		if b.status == s {
			return true
		}
	}
	return false
}
func (b *BasicWriter) BytesWritten() int {
	return b.bytes
}
func (b *BasicWriter) Tee(w io.Writer) {
	b.tee = w
}
func (b *BasicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

type FlushWriter struct {
	BasicWriter
}

func (f *FlushWriter) Flush() {
	fl := f.BasicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.Flusher = &FlushWriter{}

// HTTPFancyWriter is a HTTP writer that additionally satisfies http.CloseNotifier,
// http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type HTTPFancyWriter struct {
	BasicWriter
}

func (f *HTTPFancyWriter) CloseNotify() <-chan bool {
	cn := f.BasicWriter.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *HTTPFancyWriter) Flush() {
	fl := f.BasicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
func (f *HTTPFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.BasicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}
func (f *HTTPFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	if f.BasicWriter.tee != nil {
		n, err := io.Copy(&f.BasicWriter, r)
		f.BasicWriter.bytes += int(n)
		return n, err
	}
	rf := f.BasicWriter.ResponseWriter.(io.ReaderFrom)
	f.BasicWriter.maybeWriteHeader()
	n, err := rf.ReadFrom(r)
	f.BasicWriter.bytes += int(n)
	return n, err
}

var _ http.CloseNotifier = &HTTPFancyWriter{}
var _ http.Flusher = &HTTPFancyWriter{}
var _ http.Hijacker = &HTTPFancyWriter{}
var _ io.ReaderFrom = &HTTPFancyWriter{}

// HTTP2FancyWriter is a HTTP2 writer that additionally satisfies http.CloseNotifier,
// http.Flusher, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type HTTP2FancyWriter struct {
	BasicWriter
}

func (f *HTTP2FancyWriter) CloseNotify() <-chan bool {
	cn := f.BasicWriter.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *HTTP2FancyWriter) Flush() {
	fl := f.BasicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.CloseNotifier = &HTTP2FancyWriter{}
var _ http.Flusher = &HTTP2FancyWriter{}
