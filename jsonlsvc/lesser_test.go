package jsonlsvc

import (
	"reflect"
	"testing"
)

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
