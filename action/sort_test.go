package action

import (
	"reflect"
	"testing"
	"time"

	"github.com/yo3jones/yorg/service"
)

func TestSort(t *testing.T) {
	type test struct {
		name    string
		actions []*Action
		lessers []service.Lesser[Action]
		expect  []*Action
	}

	tests := []test{
		{
			name: "with no lessers",
			actions: []*Action{
				{Id: "foo"},
				{Id: "bar"},
			},
			lessers: []service.Lesser[Action]{},
			expect: []*Action{
				{Id: "foo"},
				{Id: "bar"},
			},
		},
		{
			name: "with all equal",
			actions: []*Action{
				{Id: "foo", Type: NextAction},
				{Id: "foo", Type: Inbox},
			},
			lessers: []service.Lesser[Action]{&OrderById{}, &OrderById{}},
			expect: []*Action{
				{Id: "foo", Type: NextAction},
				{Id: "foo", Type: Inbox},
			},
		},
		{
			name: "with less",
			actions: []*Action{
				{Id: "foo"},
				{Id: "bar"},
			},
			lessers: []service.Lesser[Action]{&OrderById{}},
			expect: []*Action{
				{Id: "bar"},
				{Id: "foo"},
			},
		},
		{
			name: "with greater",
			actions: []*Action{
				{Id: "foo"},
				{Id: "bar"},
			},
			lessers: []service.Lesser[Action]{&OrderById{desc: true}},
			expect: []*Action{
				{Id: "foo"},
				{Id: "bar"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service.Sort(tc.actions, tc.lessers)
			if !reflect.DeepEqual(tc.expect, tc.actions) {
				t.Errorf(
					"expected \n%v but got \n%v",
					tc.expect,
					tc.actions,
				)
			}
		})
	}
}

func TestLesser(t *testing.T) {
	type test struct {
		name   string
		i      *Action
		j      *Action
		l      service.Lesser[Action]
		expect int
	}

	time1, _ := time.Parse("2006-01-02", "2022-05-01")
	time2, _ := time.Parse("2006-01-02", "2022-05-02")

	tests := []test{
		{
			name:   "OrderById less",
			i:      &Action{Id: "bar"},
			j:      &Action{Id: "foo"},
			l:      &OrderById{},
			expect: -1,
		},
		{
			name:   "OrderById more",
			i:      &Action{Id: "foo"},
			j:      &Action{Id: "bar"},
			l:      &OrderById{},
			expect: 1,
		},
		{
			name:   "OrderById desc",
			i:      &Action{Id: "bar"},
			j:      &Action{Id: "foo"},
			l:      &OrderById{desc: true},
			expect: 1,
		},
		{
			name:   "OrderById equal",
			i:      &Action{Id: "foo"},
			j:      &Action{Id: "foo"},
			l:      &OrderById{},
			expect: 0,
		},
		{
			name:   "OrderByCreatedAt less",
			i:      &Action{CreatedAt: time1},
			j:      &Action{CreatedAt: time2},
			l:      &OrderByCreatedAt{},
			expect: -1,
		},
		{
			name:   "OrderByCreatedAt more",
			i:      &Action{CreatedAt: time2},
			j:      &Action{CreatedAt: time1},
			l:      &OrderByCreatedAt{},
			expect: 1,
		},
		{
			name:   "OrderByCreatedAt desc",
			i:      &Action{CreatedAt: time1},
			j:      &Action{CreatedAt: time2},
			l:      &OrderByCreatedAt{desc: true},
			expect: 1,
		},
		{
			name:   "OrderByCreatedAt equal",
			i:      &Action{CreatedAt: time1},
			j:      &Action{CreatedAt: time1},
			l:      &OrderByCreatedAt{},
			expect: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.l.Less(tc.i, tc.j)
			if got != tc.expect {
				t.Errorf("expected %d but got %d", tc.expect, got)
			}
		})
	}
}
