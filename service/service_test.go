package service

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestSpec struct {
	One string
	Two string
}

func (s *TestSpec) String() string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

type OneEquals struct {
	value string
}

func (m *OneEquals) Match(s *TestSpec) bool {
	return m.value == s.One
}

type OrderByOne struct{}

func (l *OrderByOne) Less(i, j *TestSpec) int {
	if i.One == j.One {
		return 0
	}
	if i.One < j.One {
		return -1
	}
	return 1
}

type TwoEquals struct {
	value string
}

func (m *TwoEquals) Match(s *TestSpec) bool {
	return m.value == s.Two
}

type OrderByTwo struct{}

func (l *OrderByTwo) Less(i, j *TestSpec) int {
	if i.Two == j.Two {
		return 0
	}
	if i.Two < j.Two {
		return -1
	}
	return 1
}

func TestAdd(t *testing.T) {
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

func TestMatchesAll(t *testing.T) {
	type test struct {
		name     string
		expect   bool
		spec     *TestSpec
		matchers []Matcher[TestSpec]
	}

	tests := []test{
		{
			name:   "with all match",
			expect: true,
			spec:   &TestSpec{One: "foo", Two: "bar"},
			matchers: []Matcher[TestSpec]{
				&OneEquals{"foo"},
				&TwoEquals{"bar"},
			},
		},
		{
			name:   "without match",
			expect: false,
			spec:   &TestSpec{One: "foo", Two: "bar"},
			matchers: []Matcher[TestSpec]{
				&OneEquals{"foo"},
				&TwoEquals{"baz"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MatchesAll(tc.spec, tc.matchers)
			if got != tc.expect {
				t.Errorf("expected %t but got %t", tc.expect, got)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type test struct {
		name    string
		expect  []*TestSpec
		specs   []*TestSpec
		lessers []Lesser[TestSpec]
	}

	tests := []test{
		{
			name: "with less",
			expect: []*TestSpec{
				{One: "1", Two: "1"},
				{One: "2", Two: "2"},
			},
			specs: []*TestSpec{
				{One: "1", Two: "1"},
				{One: "2", Two: "2"},
			},
			lessers: []Lesser[TestSpec]{
				&OrderByOne{},
				&OrderByTwo{},
			},
		},
		{
			name: "with more",
			expect: []*TestSpec{
				{One: "1", Two: "1"},
				{One: "2", Two: "2"},
			},
			specs: []*TestSpec{
				{One: "2", Two: "2"},
				{One: "1", Two: "1"},
			},
			lessers: []Lesser[TestSpec]{
				&OrderByOne{},
				&OrderByTwo{},
			},
		},
		{
			name: "with equal",
			expect: []*TestSpec{
				{One: "1", Two: "1"},
				{One: "1", Two: "1"},
			},
			specs: []*TestSpec{
				{One: "1", Two: "1"},
				{One: "1", Two: "1"},
			},
			lessers: []Lesser[TestSpec]{
				&OrderByOne{},
				&OrderByTwo{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Sort(tc.specs, tc.lessers)

			if !reflect.DeepEqual(tc.expect, tc.specs) {
				t.Errorf("expected \n%+v but got \n%+v", tc.expect, tc.specs)
			}
		})
	}
}
