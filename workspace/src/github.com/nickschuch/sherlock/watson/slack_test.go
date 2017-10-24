package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlackDeath(t *testing.T) {
	assert.Equal(t, "A Pod has been murdered: *foo* / *bar* / *baz*\n`sherlock inspect --s3-bucket=storage 123456`", slackMessage("foo", "bar", "baz", "storage", "123456"))
}
