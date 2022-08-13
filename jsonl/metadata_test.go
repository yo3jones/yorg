package jsonl

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/yo3jones/yorg/util"
)

type errorMetadataMarshalIndentWrapper struct{}

func (*errorMetadataMarshalIndentWrapper) marshalIndent(a any) ([]byte, error) {
	return nil, fmt.Errorf("errorMetadataMarshalIndentWrapper error")
}

func expectError(
	t *testing.T,
	expectError string,
	err error,
) (shouldContinue bool) {
	expectingError := len(expectError) > 0

	if expectingError && err == nil {
		t.Fatalf("expected an error but got nil")
		return false
	}

	if expectingError && err.Error() != expectError {
		t.Fatalf(
			"expected an error with \n[%s]\n but got \n[%s]",
			expectError,
			err.Error(),
		)
		return false
	}

	if !expectingError && err != nil {
		t.Fatal(err)
		return false
	}

	if expectingError {
		return false
	}

	return true
}

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

			if !expectError(t, tc.expectError, err) {
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

func TestSaveMetadata(t *testing.T) {
	validFilename := "valid.metadata.json"

	type test struct {
		name               string
		metadata           *JsonlMetadata
		marshalWrapper     metadataMarshalIndentWrapper
		filesExist         bool
		readOnlyPermission bool
		expectError        string
		expect             string
	}

	tests := []test{
		{
			name: "with invalid name",
			metadata: &JsonlMetadata{
				Name: "foo.jsonll",
			},
			expectError: "invalid jsonl file name foo.jsonll",
		},
		{
			name: "with marshal error",
			metadata: &JsonlMetadata{
				Name: "foo.jsonl",
			},
			marshalWrapper: &errorMetadataMarshalIndentWrapper{},
			expectError:    "errorMetadataMarshalIndentWrapper error",
		},
		{
			name: "with write error",
			metadata: &JsonlMetadata{
				Name: "valid.jsonl",
			},
			filesExist:         true,
			readOnlyPermission: true,
			expectError:        "open valid.metadata.json: permission denied",
		},
		{
			name: "with existing files",
			metadata: &JsonlMetadata{
				Name:          "valid.jsonl",
				Version:       "dev",
				MaxLineLength: 100,
			},
			filesExist: true,
			expect: `{
  "version": "dev",
  "maxLineLength": 100
}`,
		},
		{
			name: "without existing files",
			metadata: &JsonlMetadata{
				Name:          "valid.jsonl",
				Version:       "dev",
				MaxLineLength: 100,
			},
			filesExist: false,
			expect: `{
  "version": "dev",
  "maxLineLength": 100
}`,
		},
	}

	setup := func(tc test) {
		os.Remove(validFilename)
		if !tc.filesExist {
			return
		}
		var perm fs.FileMode = 0666
		if tc.readOnlyPermission {
			perm = 0444
		}
		_ = ioutil.WriteFile(
			validFilename,
			[]byte(`some data`),
			perm,
		)
	}

	teardown := func() {
		os.Remove(validFilename)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setup(tc)
			defer teardown()

			if tc.marshalWrapper != nil {
				orig := metadataMarshalIndentWrapperInst
				metadataMarshalIndentWrapperInst = tc.marshalWrapper
				defer func() { metadataMarshalIndentWrapperInst = orig }()
			}

			err := SaveMetadata(tc.metadata)

			if !expectError(t, tc.expectError, err) {
				return
			}

			var got []byte
			if got, err = ioutil.ReadFile(validFilename); err != nil {
				t.Fatal(err)
			}

			if string(got) != tc.expect {
				t.Errorf("expected \n%s\n but got \n%s", tc.expect, string(got))
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

			if !expectError(t, tc.expectError, err) {
				return
			}

			if got != tc.expect {
				t.Errorf("expected \n%s\n but got \n%s", tc.expect, got)
			}
		})
	}
}
