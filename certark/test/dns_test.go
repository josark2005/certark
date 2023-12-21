package certark_test

import (
	"encoding/json"
	"testing"

	"github.com/josark2005/certark/certark"
)

func TestGetDNSProfile(t *testing.T) {
	profile, err := certark.GetDns("cf")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("origin: ", profile)
		p, err := json.Marshal(profile)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("format: ", string(p))
		}
	}
}
