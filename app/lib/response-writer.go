package lib

import (
	"io"
	"net/http"
)

type customResponseWriter struct {
	http.ResponseWriter
	io     io.Writer
	status int
}

func (c *customResponseWriter) Write(b []byte) (int, error) {
	if c.Header().Get("Content-Type") == "" {
		c.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return c.io.Write(b)
}

func (c *customResponseWriter) WriteHeader(status int) {
	c.ResponseWriter.WriteHeader(status)
	c.status = status
}

func overrideWriter(writer http.ResponseWriter, io io.Writer) *customResponseWriter {
	return &customResponseWriter{ResponseWriter: writer, io: io, status: http.StatusOK}
}
