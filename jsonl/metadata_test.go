package jsonl

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/yo3jones/yorg/util"
)

func TestLoadMeadata(t *testing.T) {
	var err error

	invalid := `k`
	valid := `
	{
		"version": "dev",
		"maxLineLength": 1000
	}`

	os.Remove("invalid.metadata.json")
	err = ioutil.WriteFile("invalid.metadata.json", []byte(invalid), 0666)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.Remove("invalid.metadata.json")

	os.Remove("valid.metadata.json")
	err = ioutil.WriteFile("valid.metadata.json", []byte(valid), 0666)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.Remove("valid.metadata.json")

	type test struct {
		name        string
		jsonlName   string
		expect      *JsonlMetadata
		expectError string
	}

	tests := []test{
		{
			name:        "with getMetadataName error",
			jsonlName:   "test.jsonll",
			expectError: "invalid jsonl file name test.jsonll",
		},
		{
			name:        "with no file error",
			jsonlName:   "foo.jsonl",
			expectError: "open foo.metadata.json: no such file or directory",
		},
		{
			name:        "with invalid metadata",
			jsonlName:   "invalid.jsonl",
			expectError: "invalid character 'k' looking for beginning of value",
		},
		{
			name:      "with valid metadata",
			jsonlName: "valid.jsonl",
			expect: &JsonlMetadata{
				Name:          "valid.jsonl",
				Version:       "dev",
				MaxLineLength: 1000,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LoadMetadata(tc.jsonlName)

			expectingError := len(tc.expectError) > 0

			if expectingError && err == nil {
				t.Errorf("expected an error but got nil")
				return
			}

			if expectingError && err.Error() != tc.expectError {
				t.Errorf(
					"expected an error with \n%s\n but got \n%s",
					tc.expectError,
					err.Error(),
				)
				return
			}

			if !expectingError && err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf(
					"expected \n%s\n but got \n%s",
					util.Pretty(tc.expect),
					util.Pretty(got),
				)
			}
		})
	}
}

func TestGetMetadataName(t *testing.T) {
	type test struct {
		name        string
		jsonlName   string
		expect      string
		expectError string
	}

	tests := []test{
		{
			name:      "with active name",
			jsonlName: "/Users/yo3/.local/yorg/spec.jsonl",
			expect:    "/Users/yo3/.local/yorg/spec.metadata.json",
		},
		{
			name:      "with archive name",
			jsonlName: "/Users/yo3/.local/yorg/spec.archived.jsonl",
			expect:    "/Users/yo3/.local/yorg/spec.archived.metadata.json",
		},
		{
			name:        "with short name",
			jsonlName:   "1234",
			expectError: "invalid jsonl file name 1234",
		},
		{
			name:      "with invalid suffix",
			jsonlName: "/Users/yo3/.local/yorg/spec.archived.jsonll",
			expectError: strings.Join([]string{
				"invalid jsonl file name",
				"/Users/yo3/.local/yorg/spec.archived.jsonll",
			}, " "),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getMetadataName(tc.jsonlName)

			expectingError := len(tc.expectError) > 0

			if expectingError && err == nil {
				t.Errorf("expected an error but got nil")
				return
			}

			if expectingError && err.Error() != tc.expectError {
				t.Errorf(
					"expected an error with \n%s\n but got \n%s",
					tc.expectError,
					err.Error(),
				)
				return
			}

			if !expectingError && err != nil {
				t.Error(err)
				return
			}

			if got != tc.expect {
				t.Errorf("expected \n%s\n but got \n%s", tc.expect, got)
			}
		})
	}
}
