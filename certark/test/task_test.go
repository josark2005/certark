package certark_test

import (
	"testing"

	"github.com/josark2005/certark/certark"
)

func TestGetTask(t *testing.T) {
	profile, err := certark.GetTask("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestGetTaskJson(t *testing.T) {
	profile, err := certark.GetTaskJson("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestGetTaskJsonPretty(t *testing.T) {
	profile, err := certark.GetTaskJsonPretty("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestListTasks(t *testing.T) {
	users, err := certark.ListTasks()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(users)
	}
}

func TestAddTask(t *testing.T) {
	err := certark.AddAcmeUser("test", "test@test.com")
	if err != nil {
		t.Error(err)
	}
}
