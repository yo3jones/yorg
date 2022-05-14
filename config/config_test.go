package config

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

var testYaml = `
dataRoot: foo/bar
`
var expect = &Config{
	DataRoot: "foo/bar",
}

type mockUserHomeDirer struct {
	value       string
	shouldError bool
}

func (m *mockUserHomeDirer) UserHomeDir() (string, error) {
	if m.shouldError {
		return "", errors.New("mock UserHomeDir error")
	}
	return m.value, nil
}

func TestAccess(t *testing.T) {
	t.Run("with happy path", func(t *testing.T) {
		os.Remove(".yorg.yaml")
		defer os.Remove(".yorg.yaml")

		err1 := os.WriteFile(".yorg.yaml", []byte(testYaml), 0666)
		if err1 != nil {
			t.Error(err1)
		}

		a := &accessor{
			userHomeDirer: &mockUserHomeDirer{value: ""},
		}
		got, err2 := a.Access()

		if err2 != nil {
			t.Error(err1)
		}

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected \n%+v but got \n%+v", expect, got)
		}
	})

	t.Run("with UserHomeDir error", func(t *testing.T) {
		a := &accessor{
			userHomeDirer: &mockUserHomeDirer{value: "", shouldError: true},
		}

		_, err1 := a.Access()

		if err1 == nil {
			t.Fatalf("expected an error but got nil")
		}

		expect := "mock UserHomeDir error"

		if err1.Error() != expect {
			t.Errorf("expected %s but got %s", expect, err1.Error())
		}
	})

	t.Run("with Open error", func(t *testing.T) {
		os.Remove(".yorg.yaml")
		defer os.Remove(".yorg.yaml")

		a := &accessor{
			userHomeDirer: &mockUserHomeDirer{value: ""},
		}

		_, err1 := a.Access()

		if err1 == nil {
			t.Fatalf("expected an error but got nil")
		}

		expect := "open .yorg.yaml: no such file or directory"

		if err1.Error() != expect {
			t.Errorf("expected %s but got %s", expect, err1.Error())
		}
	})

	t.Run("with Decode error", func(t *testing.T) {
		os.Remove(".yorg.yaml")
		defer os.Remove(".yorg.yaml")

		err1 := os.WriteFile(".yorg.yaml", []byte(""), 0666)
		if err1 != nil {
			t.Error(err1)
		}

		a := &accessor{
			userHomeDirer: &mockUserHomeDirer{value: ""},
		}

		_, err2 := a.Access()

		if err2 == nil {
			t.Fatalf("expected an error but got nil")
		}

		expect := "EOF"

		if err2.Error() != expect {
			t.Errorf("expected %s but got %s", expect, err2.Error())
		}
	})
}

func TestUserHomeDir(t *testing.T) {
	uhd := &userHomeDirer{}
	got, err1 := uhd.UserHomeDir()

	if err1 != nil {
		t.Error(err1)
	}

	if got == "" {
		t.Errorf("expected a result but got an empty string")
	}
}

func TestLoad(t *testing.T) {
	Load()
}
