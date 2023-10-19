package certark_test

import (
	"testing"

	"github.com/jokin1999/certark/certark"
)

func TestCheckStructJsonTagExists(t *testing.T) {
	res := certark.CheckStructJsonTagExists(certark.TaskProfile{}, "test")
	if res {
		t.Error("key - test - should not exist")
	}

	res = certark.CheckStructJsonTagExists(certark.TaskProfile{}, "domain")
	if !res {
		t.Error("key - test - should exist")
	}
}
