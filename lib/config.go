package lib

import (
	"os"
	"strconv"
)

type config struct {
	Port            string // APP_PORT
	AccessLog       bool   // ACCESS_LOG
	ContentEncoding bool   // CONTENT_ENCODING
}

// Config represents its configurations
var Config *config

func init() {
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "80"
	}
	accessLog := true
	if b, err := strconv.ParseBool(os.Getenv("ACCESS_LOG")); err == nil {
		accessLog = b
	}
	contentEncoding := false
	if b, err := strconv.ParseBool(os.Getenv("CONTENT_ENCODING")); err == nil {
		contentEncoding = b
	}
	Config = &config{
		Port:            port,
		AccessLog:       accessLog,
		ContentEncoding: contentEncoding,
	}
}
