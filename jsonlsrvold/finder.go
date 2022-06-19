package jsonlsrvold

import (
	"encoding/json"
	"io"

	"github.com/yo3jones/yorg/service"
)

func (
	s *Service[S, C],
) Find(
	ms []service.Matcher[S],
	ls ...service.Lesser[S],
) ([]*S, error) {
	contexts := s.contextCreator.CreateFromMatchers(FindContextType, ms)

	specs := []*S{}
	for _, context := range contexts {
		reader, err1 := s.readCloserCreator.CreateReadCloser(
			context,
			s.homer,
		)
		if err1 != nil {
			return nil, err1
		}
		defer reader.Close()

		decoder := json.NewDecoder(reader)

		for {
			var spec S
			err2 := decoder.Decode(&spec)
			if err2 == io.EOF {
				break
			}
			if err2 != nil {
				return nil, err2
			}

			if !service.MatchesAll(&spec, ms) {
				continue
			}

			specs = append(specs, &spec)
		}
	}

	service.Sort(specs, ls)

	return specs, nil
}
