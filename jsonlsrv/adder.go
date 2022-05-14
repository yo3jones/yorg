package jsonlsrv

import (
	"encoding/json"

	"github.com/yo3jones/yorg/service"
)

func (s *Service[S, C]) Add(ms ...service.Mutator[S]) (*S, error) {
	now := s.nower.Now()

	id := s.ider.Id()
	spec := s.creator.Create(id, now)
	transaction := s.ider.Id()
	mutations := service.NewMutations(transaction, id)
	for _, m := range ms {
		m.Mutate(spec, now, mutations)
	}

	_, context := s.contextCreator.CreateFromTarget(AddContextType, spec)
	h := s.homer
	writeCloser, err1 := s.writeCloserCreator.CreateWriteCloser(context, h)
	if err1 != nil {
		return nil, err1
	}
	defer writeCloser.Close()

	encoder := json.NewEncoder(writeCloser)
	err2 := encoder.Encode(spec)
	if err2 != nil {
		return nil, err2
	}

	return spec, nil
}
