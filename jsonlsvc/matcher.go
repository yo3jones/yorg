package jsonlsvc

import "github.com/yo3jones/yorg/svc"

type Matcher[S any] interface {
	svc.Filter[S]
	Match(a *S) bool
}

func MatchesAll[T ~[]Matcher[S], S any](s *S, matchers T) bool {
	for _, m := range matchers {
		if !m.Match(s) {
			return false
		}
	}
	return true
}
