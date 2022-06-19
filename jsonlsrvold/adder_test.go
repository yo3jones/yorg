package jsonlsrvold

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var (
		now                time.Time
		writeCloserCreator *TestWriteCloserCreator
		s                  Service[TestSpec, TestContext]
		writeCloser        *TestWriteCloser
	)

	setup := func() {
		now, _ = time.Parse("2006-01-02T15:04:05", "2022-05-07T07:21:01")
		writeCloser = &TestWriteCloser{w: &strings.Builder{}}
		writeCloserCreator = &TestWriteCloserCreator{
			writeClosers: map[string]*TestWriteCloser{
				"someHome/active": writeCloser,
			},
		}

		s = Service[TestSpec, TestContext]{
			ider:               &TestIder{ids: []string{"someId"}},
			nower:              &TestNower{now},
			creator:            &TestCreator{},
			contextCreator:     &TestContextCreator{},
			homer:              &TestHomer{"someHome/"},
			writeCloserCreator: writeCloserCreator,
		}
	}

	t.Run("with happy  path", func(t *testing.T) {
		setup()

		expect := &TestSpec{
			Type: Active,
			One:  "foo",
			Two:  "bar",
			Now:  now,
		}
		expectWrite := strings.Join(
			[]string{
				strings.Join([]string{
					`{`,
					`"type":"active",`,
					`"one":"foo",`,
					`"two":"bar",`,
					`"now":"2022-05-07T07:21:01Z"`,
					`}`,
				}, ""),
				"",
			},
			"\n",
		)

		got, err1 := s.Add(
			&TypeMutator{Active},
			&OneMutator{"foo"},
			&TwoMutator{"bar"},
		)
		if err1 != nil {
			t.Error(err1)
		}

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected \n%s but got \n%s", expect, got)
		}

		if expectWrite != writeCloser.w.String() {
			t.Errorf(
				"expected \n%s but got \n%s",
				expectWrite,
				writeCloser.w.String(),
			)
		}

		if !writeCloser.closeCalled {
			t.Errorf("expected close to be called but was not")
		}
	})

	t.Run("with CreateWriteCloser error", func(t *testing.T) {
		setup()

		writeCloserCreator.shouldError = true

		_, err1 := s.Add(&OneMutator{"foo"})
		if err1 == nil {
			t.Errorf("expected an error but got nil")
		}

		expect := "mock error - CreateWriteCloser"

		if err1.Error() != expect {
			t.Errorf("expected %s but got %s", expect, err1.Error())
		}
	})

	t.Run("with Write error", func(t *testing.T) {
		setup()

		writeCloser.shouldError = true

		_, err1 := s.Add(&TypeMutator{Active}, &OneMutator{"foo"})
		if err1 == nil {
			t.Errorf("expected an error but got nil")
		}

		expect := "mock error - Write"
		if err1.Error() != expect {
			t.Errorf("expected %s but got %s", expect, err1.Error())
		}
	})
}
