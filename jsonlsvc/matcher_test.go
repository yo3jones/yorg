package jsonlsvc

import "testing"

type OneEquals struct {
	value string
}

func (m *OneEquals) Match(s *TestSpec) bool {
	return m.value == s.One
}

type TwoEquals struct {
	value string
}

func (m *TwoEquals) Match(s *TestSpec) bool {
	return m.value == s.Two
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
