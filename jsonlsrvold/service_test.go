package jsonlsrvold

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/yo3jones/yorg/config"
	"github.com/yo3jones/yorg/service"
)

type TestSpec struct {
	Type TestType  `json:"type"`
	One  string    `json:"one"`
	Two  string    `json:"two"`
	Now  time.Time `json:"now"`
}

type TestType int

const (
	Unknown TestType = iota
	Active
	Complete
)

func (t TestType) String() string {
	switch t {
	case Active:
		return "active"
	case Complete:
		return "complete"
	}
	return "unknown"
}

func FromString(t string) TestType {
	switch t {
	case "active":
		return Active
	case "complete":
		return Complete
	}
	return Unknown
}

func (t *TestType) UnmarshalJSON(b []byte) error {
	var s string

	err1 := json.Unmarshal(b, &s)
	if err1 != nil {
		return err1
	}

	*t = FromString(s)

	return nil
}

func (t TestType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TestSpec) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

type TestIder struct {
	ids []string
	i   int
}

func (ider *TestIder) Id() string {
	id := ider.ids[ider.i%len(ider.ids)]
	ider.i++
	return id
}

type TestNower struct {
	now time.Time
}

func (n *TestNower) Now() time.Time {
	return n.now
}

type TestCreator struct{}

func (tc *TestCreator) Create(id string, now time.Time) *TestSpec {
	return &TestSpec{Now: now}
}

type TypeMutator struct {
	value TestType
}

func (
	m *TypeMutator,
) Mutate(
	s *TestSpec,
	now time.Time,
	mutations *service.Mutations,
) bool {
	s.Type = m.value
	return true
}

type TypeMatcher struct {
	value TestType
}

func (m *TypeMatcher) Match(s *TestSpec) bool {
	return s.Type == m.value
}

type OneMutator struct {
	value string
}

func (
	m *OneMutator,
) Mutate(
	t *TestSpec,
	now time.Time,
	muations *service.Mutations,
) bool {
	t.One = m.value
	return true
}

type OneMatcher struct {
	value string
}

func (m *OneMatcher) Match(s *TestSpec) bool {
	return s.One == m.value
}

type OrderByOne struct {
	desc bool
}

func (o *OrderByOne) Less(i, j *TestSpec) int {
	var value int

	if i.One == j.One {
		value = 0
	} else if i.One < j.One {
		value = -1
	} else {
		value = 1
	}

	if o.desc {
		value = value * -1
	}

	return value
}

type TwoMutator struct {
	value string
}

func (
	m *TwoMutator,
) Mutate(
	t *TestSpec,
	now time.Time,
	muations *service.Mutations,
) bool {
	t.Two = m.value
	return true
}

type TwoMatcher struct {
	value string
}

func (m *TwoMatcher) Match(s *TestSpec) bool {
	return s.Two == m.value
}

type OrderByTwo struct {
	desc bool
}

func (o *OrderByTwo) Less(i, j *TestSpec) int {
	var value int

	if i.Two == j.Two {
		value = 0
	} else if i.Two < j.Two {
		value = -1
	} else {
		value = 1
	}

	if o.desc {
		value = value * -1
	}

	return value
}

type TestContextCreator struct{}

func (
	cc *TestContextCreator,
) CreateFromTarget(
	ct ContextType,
	t *TestSpec,
) (
	string,
	*TestContext,
) {
	return t.Type.String(), &TestContext{t.Type.String()}
}

func (cc *TestContextCreator) CreateFromMatchers(
	ct ContextType,
	matchers []service.Matcher[TestSpec],
) map[string]*TestContext {
	foundTypes := []TestType{}
	for _, m := range matchers {
		switch m := m.(type) {
		case *TypeMatcher:
			foundTypes = append(foundTypes, m.value)
		}
	}

	if len(foundTypes) < 1 {
		return map[string]*TestContext{
			"active":   {"active"},
			"complete": {"complete"},
		}
	}

	contexts := map[string]*TestContext{}
	for _, t := range foundTypes {
		contexts[t.String()] = &TestContext{t.String()}
	}
	return contexts
}

type TestContext struct {
	value string
}

type TestHomer struct {
	value string
}

func (t *TestHomer) Home() string {
	return t.value
}

type TestWriteCloserCreator struct {
	shouldError  bool
	writeClosers map[string]*TestWriteCloser
}

func (
	t *TestWriteCloserCreator,
) CreateWriteCloser(
	c *TestContext,
	h config.Homer,
) (
	io.WriteCloser,
	error,
) {
	if t.shouldError {
		return nil, errors.New("mock error - CreateWriteCloser")
	}

	key := fmt.Sprintf("%s%s", h.Home(), c.value)

	writeCloser, found := t.writeClosers[key]
	if !found {
		return nil, errors.New(fmt.Sprintf("no closer for %s", key))
	}

	return writeCloser, nil
}

type TestWriteCloser struct {
	w           *strings.Builder
	closeCalled bool
	shouldError bool
}

func (t *TestWriteCloser) Write(b []byte) (int, error) {
	if t.shouldError {
		return 0, errors.New("mock error - Write")
	}
	return t.w.Write(b)
}

func (t *TestWriteCloser) Close() error {
	t.closeCalled = true
	return nil
}

type TestReadCloserCreator struct {
	readers     map[string]*TestReadCloser
	shouldError bool
}

func (
	rcCreator *TestReadCloserCreator,
) CreateReadCloser(
	c *TestContext,
	h config.Homer,
) (
	io.ReadCloser,
	error,
) {
	if rcCreator.shouldError {
		return nil, errors.New("mock error - CreateReadCloser")
	}
	return rcCreator.readers[fmt.Sprintf("%s%s", h.Home(), c.value)], nil
}

type TestReadCloser struct {
	reader      io.Reader
	closed      bool
	shouldError bool
}

func (rc *TestReadCloser) Read(b []byte) (int, error) {
	if rc.shouldError {
		return 0, errors.New("mock error - Read")
	}
	return rc.reader.Read(b)
}

func (rc *TestReadCloser) Close() error {
	rc.closed = true
	return nil
}

func TestId(t *testing.T) {
	r, err1 := regexp.Compile(
		"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
	)
	if err1 != nil {
		t.Error(err1)
	}

	i := &ider{}
	got := i.Id()

	if !r.MatchString(got) {
		t.Errorf("expected a UUID format but got \n%s", got)
	}
}

func TestNow(t *testing.T) {
	n := &nower{}
	expect := time.Now()
	got := n.Now()

	if got.Unix()-expect.Unix() > 2 {
		t.Errorf("expected \n%s but got \n%s", expect, got)
	}
}

func TestNew(t *testing.T) {
	got := New[TestSpec, TestContext](
		&TestCreator{},
		&TestContextCreator{},
		&TestHomer{},
		&TestWriteCloserCreator{},
		&TestReadCloserCreator{},
	)

	if got == nil {
		t.Errorf("expected a service but got nil")
	}
}
