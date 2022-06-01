package action

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/yo3jones/yorg/config"
	"github.com/yo3jones/yorg/jsonlsrv"
	"github.com/yo3jones/yorg/service"
)

type actionCreator struct{}

func (*actionCreator) Create(id string, now time.Time) *Action {
	return &Action{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type actionContext struct {
	status      ActionStatus
	contextType jsonlsrv.ContextType
}

type actionContextCreator struct{}

func (*actionContextCreator) CreateFromTarget(
	ct jsonlsrv.ContextType,
	spec *Action,
) (string, *actionContext) {
	return spec.Status.String(), &actionContext{
		status:      spec.Status,
		contextType: ct,
	}
}

func (*actionContextCreator) CreateFromMatchers(
	ct jsonlsrv.ContextType,
	matchers []service.Matcher[Action],
) map[string]*actionContext {
	contexts := map[string]*actionContext{}

	for _, matcher := range matchers {
		switch matcher := matcher.(type) {
		case *StatusEquals:
			contexts[matcher.value.String()] = &actionContext{
				status:      matcher.value,
				contextType: ct,
			}
		}
	}

	if len(contexts) <= 0 {
		contexts[Active.String()] = &actionContext{
			status:      Active,
			contextType: ct,
		}
		contexts[Complete.String()] = &actionContext{
			status:      Complete,
			contextType: ct,
		}
	}

	return contexts
}

func buildFilePath(c *actionContext, h config.Homer) string {
	fileName := "actions"
	if c.status == Complete {
		fileName = fmt.Sprintf("%s-complete", fileName)
	}
	return filepath.Join(h.Home(), fmt.Sprintf("%s.jsonl", fileName))
}

type actionWriteCloserCreator struct{}

func (*actionWriteCloserCreator) CreateWriteCloser(
	c *actionContext,
	h config.Homer,
) (io.WriteCloser, error) {
	return os.OpenFile(
		buildFilePath(c, h),
		os.O_RDWR,
		0644,
	)
}

type actionReadCloserCreator struct{}

func (*actionReadCloserCreator) CreateReadCloser(
	c *actionContext,
	h config.Homer,
) (io.ReadCloser, error) {
	return os.Open(buildFilePath(c, h))
}

func NewService(homer config.Homer) service.Service[Action] {
	return jsonlsrv.New[Action, actionContext](
		&actionCreator{},
		&actionContextCreator{},
		homer,
		&actionWriteCloserCreator{},
		&actionReadCloserCreator{},
	)
}
