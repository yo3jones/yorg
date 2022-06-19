package svc

import (
	"time"
)

type Service[S Spec, F Filter[S], O Order[S]] interface {
	Adder[S]
	// Finder[S, F, O]
}

type Spec interface {
	Ider
}

type Ider interface {
	Id() string
	SetId(id string)
}

type Adder[S Spec] interface {
	Add(ms ...Mutator[S]) (*S, error)
}

type Finder[S any, F Filter[S], O Order[S]] interface {
	Find(filters []F, orders ...O) ([]*S, error)
}

type Mutator[S Spec] interface {
	Mutate(s *S, t time.Time, m *Mutations) bool
}

type Filter[S any] interface{}

type Order[S any] interface{}
