package certark_test

import (
	"testing"

	"github.com/josark2005/certark/certark"
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
