package jsonlsvc

import (
	"sort"

	"github.com/yo3jones/yorg/svc"
)

type Lesser[S any] interface {
	svc.Filter[S]
	Less(i, j *S) int
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
