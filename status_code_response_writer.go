// Created by davidterranova on 16/10/2017.

package hserver

import "net/http"

type StatusCodeResponseWriter struct {
	status int
	http.ResponseWriter
	wroteHeader bool
}

func NewStatusCodeResponseWriter(w http.ResponseWriter) *StatusCodeResponseWriter {
	return &StatusCodeResponseWriter{
		http.StatusOK,
		w,
		false,
	}
}

func (w *StatusCodeResponseWriter) Status() int {
	return w.status
}

func (w *StatusCodeResponseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(w.status)
	}
	return w.ResponseWriter.Write(data)
}

func (w *StatusCodeResponseWriter) WriteHeader(code int) {
	// Write the status code onward.
	w.ResponseWriter.WriteHeader(code)

	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}
