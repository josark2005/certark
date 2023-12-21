package certark_test

import (
	"testing"

	"github.com/josark2005/certark/certark"
)

func TestGetAcmeUser(t *testing.T) {
	profile, err := certark.GetAcmeUser("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestGetAcmeUserJson(t *testing.T) {
	profile, err := certark.GetAcmeUserJson("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestGetAcmeUserJsonPretty(t *testing.T) {
	profile, err := certark.GetAcmeUserJsonPretty("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestListAcmeUsers(t *testing.T) {
	users, err := certark.ListAcmeUsers()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(users)
	}
}

func TestAddAcmeUser(t *testing.T) {
	err := certark.AddAcmeUser("test", "test@test.com")
	if err != nil {
		t.Error(err)
	}
}
