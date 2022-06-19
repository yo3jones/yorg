package svc

import (
	"reflect"
	"testing"
)

func TestNewMutations(t *testing.T) {
	expect := &Mutations{
		Transaction: "someTransaction",
		Id:          "someId",
		From:        make(map[string]any),
		To:          make(map[string]any),
	}
	got := NewMutations("someTransaction", "someId")

	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected \n%+v but got \n %+v", expect, got)
	}
}

func TestMutationsAdd(t *testing.T) {
	t.Run("without exists", func(t *testing.T) {
		expect := &Mutations{
			Transaction: "someTransaction",
			Id:          "someId",
			From:        map[string]any{"one": "foo"},
			To:          map[string]any{"one": "bar"},
		}
		got := &Mutations{
			Transaction: "someTransaction",
			Id:          "someId",
			From:        make(map[string]any),
			To:          make(map[string]any),
		}

		got.Add("one", "foo", "bar")

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected \n%+v but got \n %+v", expect, got)
		}
	})

	t.Run("with exists", func(t *testing.T) {
		expect := &Mutations{
			Transaction: "someTransaction",
			Id:          "someId",
			From:        map[string]any{"one": "foo"},
			To:          map[string]any{"one": "fiz"},
		}
		got := &Mutations{
			Transaction: "someTransaction",
			Id:          "someId",
			From:        make(map[string]any),
			To:          make(map[string]any),
		}

		got.Add("one", "foo", "bar")
		got.Add("one", "bar", "fiz")

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected \n%+v but got \n %+v", expect, got)
		}
	})
}
