package lib

import (
	"os"
	"strconv"
)

type config struct {
	AccessLog       bool // ACCESS_LOG
	ContentEncoding bool // CONTENT_ENCODING
}

// Config represents its configurations
var Config *config

func init() {
	setup()
}

func setup() {
	accessLog := true
	if b, err := strconv.ParseBool(os.Getenv("ACCESS_LOG")); err == nil {
		accessLog = b
	}
	contentEncoding := true
	if b, err := strconv.ParseBool(os.Getenv("CONTENT_ENCODING")); err == nil {
		contentEncoding = b
	}
	Config = &config{
		AccessLog:       accessLog,
		ContentEncoding: contentEncoding,
	}
}
