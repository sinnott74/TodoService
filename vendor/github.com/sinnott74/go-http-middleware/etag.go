package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"strconv"
)

// DefaultEtag middleware which uses MD5 as its hashing function
func DefaultEtag(next http.Handler) http.Handler {
	return Etag(md5.New)(next)
}

// Etag middleware which handles adding an ETag header to the response
// An ETag is a hash of a resource that client's/browser use to cache resourse that are unchanged.
// It allows the server to skip sending the resource over the object if the client has it already
// A StatusNotModified (304) is returned when the client's resource is up to date.
// Client's set the If-None-Match header to send their cached ETag for a resource
func Etag(newHash func() hash.Hash) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			hash := newHash()
			etagWriter := &etagWriter{rw: w, hash: hash, buf: bytes.NewBuffer(nil)}
			next.ServeHTTP(etagWriter, r)

			if !isStatusOk(etagWriter.status) || etagWriter.status == http.StatusNoContent || etagWriter.buf.Len() == 0 {
				etagWriter.writeResponse()
				return
			}

			reqEtag := r.Header.Get("If-None-Match")
			responseEtag := etagWriter.etag()
			w.Header().Set("Etag", responseEtag)

			if responseEtag == reqEtag {
				w.WriteHeader(http.StatusNotModified)
				w.Write(nil)
			} else {
				etagWriter.writeResponse()
			}
		})
	}
}

// etagWriter is an stuct which implements the ResponseWriter interface
// Its responsible for capturing whats written the response & hashing it
// so that it can be used as an etag header
type etagWriter struct {
	rw     http.ResponseWriter
	hash   hash.Hash
	buf    *bytes.Buffer
	status int
}

// Header delegates to the http response Header
func (w etagWriter) Header() http.Header {
	return w.rw.Header()
}

// WriteHeader sets the status of this writer to be set in the http response later
func (w *etagWriter) WriteHeader(status int) {
	w.status = status
}

// Write the bytes to both the buffer & the hash
func (w *etagWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.buf.Write(b)
	l, err := w.hash.Write(b)
	return l, err
}

// writeResponse writes the buffer to the response
func (w *etagWriter) writeResponse() {
	w.rw.WriteHeader(w.status)
	w.rw.Write(w.buf.Bytes())
}

// sumHash finishes & returns the hashed response
func (w *etagWriter) sumHash() []byte {
	return w.hash.Sum(nil)
}

// etag outputs etag for the response, which contains the hash response
func (w *etagWriter) etag() string {
	sumHash := w.sumHash()
	base64Hash := base64.StdEncoding.EncodeToString(sumHash)
	len := strconv.FormatInt(int64(w.buf.Len()), 16) // hexidecimal
	return fmt.Sprintf("W/\"%v-%v\"", len, base64Hash)
}

// isStatusOk check is the given http status is in the 2xx range
func isStatusOk(status int) bool {
	return status >= 200 && status < 300
}
