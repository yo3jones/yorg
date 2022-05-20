package action

import (
	"testing"

	"github.com/yo3jones/yorg/service"
)

func TestQueries(t *testing.T) {
	type test struct {
		name        string
		matcher     service.Matcher[Action]
		action      *Action
		expect      bool
	}

	tests := []test{
		{
			name:        "IdEquals with match",
			matcher:     &IdEquals{"foo"},
			action:      &Action{Id: "foo"},
			expect:      true,
		},
		{
			name:        "IdEquals without match",
			matcher:     &IdEquals{"bar"},
			action:      &Action{Id: "foo"},
			expect:      false,
		},
		{
			name:        "TypeEquals with match",
			matcher:     &TypeEquals{NextAction},
			action:      &Action{Type: NextAction},
			expect:      true,
		},
		{
			name:        "TypeEquals without match",
			matcher:     &TypeEquals{Inbox},
			action:      &Action{Type: NextAction},
			expect:      false,
		},
		{
			name:        "StatusEquals with match",
			matcher:     &StatusEquals{Complete},
			action:      &Action{Status: Complete},
			expect:      true,
		},
		{
			name:        "StatusEquals without match",
			matcher:     &StatusEquals{Active},
			action:      &Action{Status: Complete},
			expect:      false,
		},
		{
			name:        "ProjectEquals with match",
			matcher:     &ProjectEquals{"foo"},
			action:      &Action{Project: "foo"},
			expect:      true,
		},
		{
			name:        "ProjectEquals without match",
			matcher:     &ProjectEquals{"bar"},
			action:      &Action{Project: "foo"},
			expect:      false,
		},
		{
			name:        "HasTag with match",
			matcher:     &HasTag{"Foo"},
			action:      &Action{Tags: []string{"bar", "foo"}},
			expect:      true,
		},
		{
			name:        "HasTag without match",
			matcher:     &HasTag{"fiz"},
			action:      &Action{Tags: []string{"bar", "foo"}},
			expect:      false,
		},
		{
			name:        "DoesNotHaveTag with match",
			matcher:     &DoesNotHaveTag{"fiz"},
			action:      &Action{Tags: []string{"bar", "foo"}},
			expect:      true,
		},
		{
			name:        "DoesNotHaveTag without match",
			matcher:     &DoesNotHaveTag{"Foo"},
			action:      &Action{Tags: []string{"bar", "foo"}},
			expect:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.matcher.Match(tc.action)
			if tc.expect != got {
				t.Errorf("expected %t but got %t", tc.expect, got)
			}
		})
	}
}
