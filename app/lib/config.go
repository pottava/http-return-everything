package lib

import (
	"os"
	"strconv"
)

type config struct {
	AccessLog       bool // ACCESS_LOG
	AccessDetailLog bool // ACCESS_DETAIL_LOG
	ContentEncoding bool // CONTENT_ENCODING
	EnabledAWS      bool // ENABLE_AWS
	EnabledGCP      bool // ENABLE_GCP
}

var (
	ver    = "dev"
	commit string
	date   string
)

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
	accessDetailLog := false
	if b, err := strconv.ParseBool(os.Getenv("ACCESS_DETAIL_LOG")); err == nil {
		accessDetailLog = b
	}
	contentEncoding := true
	if b, err := strconv.ParseBool(os.Getenv("CONTENT_ENCODING")); err == nil {
		contentEncoding = b
	}
	enableAWS := true
	if b, err := strconv.ParseBool(os.Getenv("ENABLE_AWS")); err == nil {
		enableAWS = b
	}
	enableGCP := true
	if b, err := strconv.ParseBool(os.Getenv("ENABLE_GCP")); err == nil {
		enableGCP = b
	}
	Config = &config{
		AccessLog:       accessLog,
		AccessDetailLog: accessDetailLog,
		ContentEncoding: contentEncoding,
		EnabledAWS:      enableAWS,
		EnabledGCP:      enableGCP,
	}
}
