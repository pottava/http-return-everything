package lib

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
)

// Wrap wraps HTTP request handler
func Wrap(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case eqauls(r, "/health"):
			w.WriteHeader(http.StatusOK)

		case eqauls(r, "/version"):
			if len(commit) > 0 && len(date) > 0 {
				fmt.Fprintf(w, "%s-%s (built at %s)\n", ver, commit, date)
				return
			}
			fmt.Fprintln(w, ver)

		default:
			proc := time.Now()
			addr := r.RemoteAddr
			if ip, found := header(r, "X-Forwarded-For"); found {
				addr = ip[0]
			}
			ioWriter := w.(io.Writer)
			if encodings, found := header(r, "Accept-Encoding"); found && Config.ContentEncoding {
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
			if Config.CORSOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", Config.CORSOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,PATCH,HEAD")
				w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
				w.Header().Set("Access-Control-Expose-Headers", "Content-Type,Authorization,Date")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}
			writer := overrideWriter(w, ioWriter)
			handler.ServeHTTP(writer, r)

			if Config.AccessLog {
				if Config.AccessDetailLog {
					marshaled, _ := json.Marshal(models.HTTPRequest{
						Protocol:   swag.String(r.Proto),
						Method:     swag.String(r.Method),
						Host:       swag.String(r.Host),
						RemoteAddr: swag.String(r.RemoteAddr),
						URI:        swag.String(r.RequestURI),
						Headers:    r.Header,
						Form:       r.Form,
						PostForm:   r.PostForm,
					})
					log.Printf("[%s] %.3f %d %s",
						addr, time.Since(proc).Seconds(),
						writer.status, marshaled)
				} else {
					log.Printf("[%s] %.3f %d %s %s",
						addr, time.Since(proc).Seconds(),
						writer.status, r.Method, r.URL)
				}
			}
		}
	})
}

func eqauls(r *http.Request, url string) bool {
	return url == r.URL.Path
}

func header(r *http.Request, key string) (values []string, found bool) {
	if r.Header == nil {
		return
	}
	for k, v := range r.Header {
		if strings.EqualFold(k, key) && len(v) > 0 {
			return v, true
		}
	}
	return
}

func splitCsvLine(data string) []string {
	splitted := strings.Split(data, ",")
	parsed := make([]string, len(splitted))
	for i, val := range splitted {
		parsed[i] = strings.TrimSpace(val)
	}
	return parsed
}
