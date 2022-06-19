package jsonlsrvold

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/yo3jones/yorg/service"
)

func TestFind(t *testing.T) {
	var (
		s                 *Service[TestSpec, TestContext]
		readCloserCreator *TestReadCloserCreator
		readClosers       map[string]*TestReadCloser
	)
	now, _ := time.Parse("2006-01-02T15:04:05", "2022-05-07T12:24:57")

	setup := func() {
		readClosers = map[string]*TestReadCloser{
			"someHome/active":   {},
			"someHome/complete": {},
		}
		readCloserCreator = &TestReadCloserCreator{readers: readClosers}
		s = &Service[TestSpec, TestContext]{
			homer:             &TestHomer{"someHome/"},
			contextCreator:    &TestContextCreator{},
			readCloserCreator: readCloserCreator,
		}
	}

	t.Run("when succeeds", func(t *testing.T) {
		type test struct {
			name     string
			expect   []*TestSpec
			closed   []string
			inputs   map[string][][]string
			matchers []service.Matcher[TestSpec]
			lessers  []service.Lesser[TestSpec]
		}

		tests := []test{
			{
				name: "with single reader",
				expect: []*TestSpec{
					{
						Type: Active,
						One:  "foo",
						Two:  "bar",
						Now:  now,
					},
				},
				closed: []string{"someHome/active"},
				inputs: map[string][][]string{
					"someHome/active": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
					"someHome/complete": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
				},
				matchers: []service.Matcher[TestSpec]{
					&TypeMatcher{Active},
				},
			},
			{
				name: "with filtering",
				expect: []*TestSpec{
					{
						Type: Active,
						One:  "foo",
						Two:  "bar",
						Now:  now,
					},
				},
				closed: []string{"someHome/active"},
				inputs: map[string][][]string{
					"someHome/active": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{
							`{`,
							`"type":"active",`,
							`"one":"fiz",`,
							`"two":"baz",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
					"someHome/complete": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
				},
				matchers: []service.Matcher[TestSpec]{
					&TypeMatcher{Active},
					&OneMatcher{"foo"},
					&TwoMatcher{"bar"},
				},
			},
			{
				name: "with multiple readers",
				expect: []*TestSpec{
					{
						Type: Active,
						One:  "foo",
						Two:  "bar",
						Now:  now,
					},
					{
						Type: Active,
						One:  "foo",
						Two:  "bar",
						Now:  now,
					},
				},
				closed: []string{"someHome/active", "someHome/complete"},
				inputs: map[string][][]string{
					"someHome/active": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
					"someHome/complete": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
				},
				matchers: []service.Matcher[TestSpec]{
					&OneMatcher{"foo"},
					&TwoMatcher{"bar"},
				},
			},
			{
				name: "with sorting",
				expect: []*TestSpec{
					{
						Type: Active,
						One:  "foo",
						Two:  "bar",
						Now:  now,
					},
					{
						Type: Active,
						One:  "foo",
						Two:  "baz",
						Now:  now,
					},
				},
				closed: []string{"someHome/active"},
				inputs: map[string][][]string{
					"someHome/active": {
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"baz",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{
							`{`,
							`"type":"active",`,
							`"one":"foo",`,
							`"two":"bar",`,
							`"now":"2022-05-07T12:24:57Z"`,
							`}`,
						},
						{},
					},
				},
				matchers: []service.Matcher[TestSpec]{
					&TypeMatcher{Active},
				},
				lessers: []service.Lesser[TestSpec]{
					&OrderByOne{},
					&OrderByTwo{},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				setup()
				for key, jsons := range tc.inputs {
					lines := []string{}
					for _, json := range jsons {
						lines = append(lines, strings.Join(json, ""))
					}
					readClosers[key].reader = strings.NewReader(
						strings.Join(lines, "\n"),
					)
				}

				got, err1 := s.Find(tc.matchers, tc.lessers...)

				if err1 != nil {
					t.Error(err1)
				}

				for _, key := range tc.closed {
					if !readClosers[key].closed {
						t.Errorf("expected %s to be closed but was not", key)
					}
				}

				if !reflect.DeepEqual(got, tc.expect) {
					t.Errorf("expected %v but got %v", tc.expect, got)
				}
			})
		}
	})

	t.Run("when fails", func(t *testing.T) {
		t.Run("with CreateReadCloser error", func(t *testing.T) {
			setup()
			readCloserCreator.shouldError = true

			_, err1 := s.Find([]service.Matcher[TestSpec]{
				&TypeMatcher{Active},
				&OneMatcher{"foo"},
				&TwoMatcher{"bar"},
			})

			if err1 == nil {
				t.Errorf("expected an error but got nil")
			}

			expect := "mock error - CreateReadCloser"
			if err1.Error() != expect {
				t.Errorf("expected %s but got %s", expect, err1.Error())
			}
		})

		t.Run("with Read error", func(t *testing.T) {
			setup()
			readClosers["someHome/active"].shouldError = true

			_, err1 := s.Find([]service.Matcher[TestSpec]{
				&TypeMatcher{Active},
				&OneMatcher{"foo"},
				&TwoMatcher{"bar"},
			})

			if err1 == nil {
				t.Errorf("expected an error but got nil")
			}

			expect := "mock error - Read"
			if err1.Error() != expect {
				t.Errorf("expected %s but got %s", expect, err1.Error())
			}
		})
	})
}
