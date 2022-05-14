package jsonlsrv

import (
	"io"
	"time"

	"github.com/yo3jones/yorg/config"
	"github.com/yo3jones/yorg/service"
)

type Spec interface{}

type Ider interface {
	Id() string
}

type Nower interface {
	Now() time.Time
}

type Creator[S Spec] interface {
	Create(id string, now time.Time) *S
}

type ContextCreator[S Spec, C any] interface {
	CreateFromTarget(ct ContextType, spec *S) (string, *C)
	CreateFromMatchers(
		ct ContextType,
		matchers []service.Matcher[S],
	) map[string]*C
}

type ContextType int

const (
	UnknownContextType ContextType = iota
	AddContextType
	FindContextType
)

type WriteCloserCreator[C any] interface {
	CreateWriteCloser(c *C, h config.Homer) (io.WriteCloser, error)
}

type ReadCloserCreator[C any] interface {
	CreateReadCloser(c *C, h config.Homer) (io.ReadCloser, error)
}

type Service[S Spec, C any] struct {
	ider               Ider
	nower              Nower
	creator            Creator[S]
	contextCreator     ContextCreator[S, C]
	homer              config.Homer
	writeCloserCreator WriteCloserCreator[C]
	readCloserCreator  ReadCloserCreator[C]
}
