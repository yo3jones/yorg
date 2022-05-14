package service

import (
	"sort"
	"time"
)

type Service[S any] interface {
	Adder[S]
	Finder[S]
}

type Adder[S any] interface {
	Add(ms ...Mutator[S]) (*S, error)
}

type Finder[S any] interface {
	Find(ms []Matcher[S], ls ...Lesser[S]) ([]*S, error)
}

type Mutator[S any] interface {
	Mutate(s *S, t time.Time, m *Mutations) bool
}

type Mutations struct {
	Transaction string
	Id          string
	From        map[string]any
	To          map[string]any
}

func (m *Mutations) Add(field string, from any, to any) {
	_, exists := m.From[field]

	if !exists {
		m.From[field] = from
	}
	m.To[field] = to
}

func NewMutations(transaction string, id string) *Mutations {
	return &Mutations{
		Transaction: transaction,
		Id:          id,
		From:        make(map[string]any),
		To:          make(map[string]any),
	}
}

type Matcher[S any] interface {
	Match(a *S) bool
}

type Lesser[S any] interface {
	Less(i, j *S) int
}

func MatchesAll[T ~[]Matcher[S], S any](s *S, matchers T) bool {
	for _, m := range matchers {
		if !m.Match(s) {
			return false
		}
	}
	return true
}

func Sort[T ~[]Lesser[S], I ~[]*S, S any](ss I, lessers T) {
	sort.SliceStable(ss, func(i, j int) bool {
		for _, lesser := range lessers {
			r := lesser.Less(ss[i], ss[j])
			if r == 0 {
				continue
			}
			if r < 0 {
				return true
			} else {
				return false
			}
		}
		return false
	})
}
