package jsonl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type JsonlMetadata struct {
	Name          string `json:"-"`
	Version       string `json:"version"`
	MaxLineLength int    `json:"maxLineLength"`
}

func LoadMetadata(name string) (metadata *JsonlMetadata, err error) {
	var (
		metadataName string
		b            []byte
	)

	if metadataName, err = getMetadataName(name); err != nil {
		return nil, err
	}

	if b, err = ioutil.ReadFile(metadataName); err != nil {
		return nil, err
	}

	metadata = &JsonlMetadata{}
	if err = json.Unmarshal(b, metadata); err != nil {
		return nil, err
	}

	metadata.Name = name

	return metadata, nil
}

func SaveMetadata(metadata *JsonlMetadata) (err error) {
	var (
		filename string
		data     []byte
	)
	if filename, err = getMetadataName(metadata.Name); err != nil {
		return err
	}

	data, err = metadataMarshalIndentWrapperInst.marshalIndent(metadata)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		return err
	}

	return nil
}

func getMetadataName(jsonlName string) (string, error) {
	jsonlNameLen := len(jsonlName)

	if jsonlNameLen <= 6 || !strings.HasSuffix(jsonlName, ".jsonl") {
		return "", fmt.Errorf("invalid jsonl file name %s", jsonlName)
	}

	suffixLen := len(".jsonl")

	name := fmt.Sprintf(
		"%s.metadata.json",
		jsonlName[:jsonlNameLen-suffixLen],
	)

	return name, nil
}

var metadataMarshalIndentWrapperInst metadataMarshalIndentWrapper = &metadataMarshalIndentWrapperImpl{}

type metadataMarshalIndentWrapper interface {
	marshalIndent(a any) ([]byte, error)
}

type metadataMarshalIndentWrapperImpl struct{}

func (*metadataMarshalIndentWrapperImpl) marshalIndent(a any) ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}
