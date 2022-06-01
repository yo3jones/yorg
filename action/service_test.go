package action

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/yo3jones/yorg/jsonlsrv"
	"github.com/yo3jones/yorg/service"
)

func TestNewService(t *testing.T) {
	got := NewService(nil)

	if got == nil {
		t.Errorf("expected a new service but got nil")
	}
}

func TestCreate(t *testing.T) {
	now := time.Now()
	id := "some id"

	ac := &actionCreator{}
	got := ac.Create(id, now)

	if got == nil {
		t.Errorf("expected an action to have been created  but got nil")
	}

	if got.Id != id {
		t.Errorf(
			"expected Id to have been set with [%s] but got [%s]",
			id,
			got.Id,
		)
	}

	if got.CreatedAt != now {
		t.Errorf(
			"expected CreatedAt to have been set with \n%s but got \n%s",
			now,
			got.CreatedAt,
		)
	}

	if got.UpdatedAt != now {
		t.Errorf(
			"expected UpdatedAt to have been set with \n%s but got \n%s",
			now,
			got.UpdatedAt,
		)
	}
}

func TestCreateFromTarget(t *testing.T) {
	type test struct {
		name          string
		contextType   jsonlsrv.ContextType
		action        *Action
		expectKey     string
		expectContext *actionContext
	}

	tests := []test{
		{
			name:        "active",
			contextType: jsonlsrv.AddContextType,
			action:      &Action{Status: Active},
			expectKey:   "active",
			expectContext: &actionContext{
				status:      Active,
				contextType: jsonlsrv.AddContextType,
			},
		},
		{
			name:        "complete",
			contextType: jsonlsrv.FindContextType,
			action:      &Action{Status: Complete},
			expectKey:   "complete",
			expectContext: &actionContext{
				status:      Complete,
				contextType: jsonlsrv.FindContextType,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			contextCreator := &actionContextCreator{}
			gotKey, gotContext := contextCreator.CreateFromTarget(
				tc.contextType,
				tc.action,
			)

			if gotKey != tc.expectKey {
				t.Errorf("expected %s but got %s", tc.expectKey, gotKey)
			}

			if !reflect.DeepEqual(gotContext, tc.expectContext) {
				t.Errorf(
					"expected \n%+v but got \n%+v",
					tc.expectContext,
					gotContext,
				)
			}
		})
	}
}

func TestCreateFromMatchers(t *testing.T) {
	type test struct {
		name        string
		contextType jsonlsrv.ContextType
		matchers    []service.Matcher[Action]
		expect      map[string]*actionContext
	}

	tests := []test{
		{
			name:        "with only active",
			contextType: jsonlsrv.AddContextType,
			matchers: []service.Matcher[Action]{
				&StatusEquals{Active},
				&StatusEquals{Active},
			},
			expect: map[string]*actionContext{
				"active": {
					status:      Active,
					contextType: jsonlsrv.AddContextType,
				},
			},
		},
		{
			name:        "with only complete",
			contextType: jsonlsrv.FindContextType,
			matchers:    []service.Matcher[Action]{&StatusEquals{Complete}},
			expect: map[string]*actionContext{
				"complete": {
					status:      Complete,
					contextType: jsonlsrv.FindContextType,
				},
			},
		},
		{
			name:        "with both",
			contextType: jsonlsrv.FindContextType,
			matchers: []service.Matcher[Action]{
				&StatusEquals{Active},
				&StatusEquals{Complete},
			},
			expect: map[string]*actionContext{
				"active": {
					status:      Active,
					contextType: jsonlsrv.FindContextType,
				},
				"complete": {
					status:      Complete,
					contextType: jsonlsrv.FindContextType,
				},
			},
		},
		{
			name:        "with none",
			contextType: jsonlsrv.FindContextType,
			matchers:    []service.Matcher[Action]{},
			expect: map[string]*actionContext{
				"active": {
					status:      Active,
					contextType: jsonlsrv.FindContextType,
				},
				"complete": {
					status:      Complete,
					contextType: jsonlsrv.FindContextType,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			contextCreator := &actionContextCreator{}
			got := contextCreator.CreateFromMatchers(
				tc.contextType,
				tc.matchers,
			)

			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("expected \n%#v but got \n%#v", tc.expect, got)
			}
		})
	}
}

type TestHomer struct {
	called bool
}

func (h *TestHomer) Home() string {
	h.called = true
	return ""
}

func TestCreateWriteCloser(t *testing.T) {
	t.Run("with active", func(t *testing.T) {
		os.Create("actions.jsonl")
		defer os.Remove("actions.jsonl")

		writeCloseCreator := &actionWriteCloserCreator{}
		h := &TestHomer{}
		got, err1 := writeCloseCreator.CreateWriteCloser(
			&actionContext{
				status: Active,
			},
			h,
		)

		if err1 != nil {
			t.Fatal(err1)
		}

		if got == nil {
			t.Errorf("got nil")
		}

		if !h.called {
			t.Error("expected Homer.Home() to be called but was not")
		}
	})

	t.Run("with complete", func(t *testing.T) {
		os.Create("actions-complete.jsonl")
		defer os.Remove("actions-complete.jsonl")

		writeCloseCreator := &actionWriteCloserCreator{}
		h := &TestHomer{}
		got, err1 := writeCloseCreator.CreateWriteCloser(
			&actionContext{
				status: Complete,
			},
			h,
		)
		if err1 != nil {
			t.Fatal(err1)
		}

		if got == nil {
			t.Errorf("got nil")
		}

		if !h.called {
			t.Error("expected Homer.Home() to be called but was not")
		}
	})
}

func TestCreateReadCloser(t *testing.T) {
	os.Create("actions.jsonl")
	defer os.Remove("actions.jsonl")

	readCloserCreator := &actionReadCloserCreator{}
	h := &TestHomer{}
	got, err1 := readCloserCreator.CreateReadCloser(
		&actionContext{
			status: Active,
		},
		h,
	)

	if err1 != nil {
		t.Fatal(err1)
	}

	if got == nil {
		t.Errorf("got nil")
	}

	if !h.called {
		t.Error("expected Homer.Home() to be called but was not")
	}
}
