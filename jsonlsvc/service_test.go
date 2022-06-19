package jsonlsvc

import "encoding/json"

type TestSpec struct {
	One string
	Two string
}

func (s *TestSpec) String() string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}
