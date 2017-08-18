package lib

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Wrap wraps HTTP request handler
func Wrap(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proc := time.Now()
		addr := r.RemoteAddr
		if ip, found := Header(r, "X-Forwarded-For"); found {
			addr = ip[0]
		}
		ioWriter := w.(io.Writer)
		if encodings, found := Header(r, "Accept-Encoding"); found && Config.ContentEncoding {
			for _, encoding := range splitCsvLine(encodings[0]) {
				if encoding == "gzip" {
					w.Header().Set("Content-Encoding", "gzip")
					g := gzip.NewWriter(w)
					defer g.Close()
					ioWriter = g
					break
				}
				if encoding == "deflate" {
					w.Header().Set("Content-Encoding", "deflate")
					z := zlib.NewWriter(w)
					defer z.Close()
					ioWriter = z
					break
				}
			}
		}
		writer := overrideWriter(w, ioWriter)
		f(writer, r)

		if Config.AccessLog {
			log.Printf("[%s] %.3f %d %s %s",
				addr, time.Now().Sub(proc).Seconds(),
				writer.status, r.Method, r.URL)
		}
	})
}

func splitCsvLine(data string) []string {
	splitted := strings.SplitN(data, ",", -1)
	parsed := make([]string, len(splitted))
	for i, val := range splitted {
		parsed[i] = strings.TrimSpace(val)
	}
	return parsed
}
