package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	os.Setenv("ACCESS_LOG", "false")
	os.Setenv("ACCESS_DETAIL_LOG", "true")
	os.Setenv("CONTENT_ENCODING", "false")

	setup()

	assert.Equal(t, false, Config.AccessLog, "Expected an error when AccessLog is not false")
	assert.Equal(t, true, Config.AccessDetailLog, "Expected an error when AccessDetailLog is not true")
	assert.Equal(t, false, Config.ContentEncoding, "Expected an error when ContentEncoding is not false")
}
