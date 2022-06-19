package jsonlsvc

import (
	"github.com/yo3jones/yorg/svc"
)

type JsonlSvc[S svc.Spec] interface {
	svc.Service[S, Matcher[S], Lesser[S]]
}

type jsonlSvc[S svc.Spec] struct{}

func New[S svc.Spec]() JsonlSvc[S] {
	return &jsonlSvc[S]{}
}
