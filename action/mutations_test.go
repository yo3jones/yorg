package action

import (
	"reflect"
	"testing"
	"time"

	"github.com/yo3jones/yorg/service"
)

func TestMutate(t *testing.T) {
	type test struct {
		name              string
		mutations         func() *service.Mutations
		mutator           service.Mutator[Action]
		in                Action
		after             Action
		expect            bool
		expectedMutations *service.Mutations
	}

	now := time.Now()

	tests := []test{
		{
			"type with change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateType{NextAction},
			Action{Type: Inbox},
			Action{Type: NextAction},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From:        map[string]any{"type": Inbox},
				To:          map[string]any{"type": NextAction},
			},
		},
		{
			"type without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateType{NextAction},
			Action{Type: NextAction},
			Action{Type: NextAction},
			false,
			service.NewMutations("", ""),
		},
		{
			"action with change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateAction{"bar"},
			Action{Action: "foo"},
			Action{Action: "bar"},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From:        map[string]any{"action": "foo"},
				To:          map[string]any{"action": "bar"},
			},
		},
		{
			"action without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateAction{"foo"},
			Action{Action: "foo"},
			Action{Action: "foo"},
			false,
			service.NewMutations("", ""),
		},
		{
			"status active to complete",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateStatus{Complete},
			Action{Status: Active},
			Action{Status: Complete, CompletedAt: now},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From: map[string]any{
					"status":      Active,
					"completedAt": time.Time{},
				},
				To: map[string]any{
					"status":      Complete,
					"completedAt": now,
				},
			},
		},
		{
			"status complete to active",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateStatus{Active},
			Action{Status: Complete, CompletedAt: now},
			Action{Status: Active},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From: map[string]any{
					"status":      Complete,
					"completedAt": now,
				},
				To: map[string]any{
					"status":      Active,
					"completedAt": time.Time{},
				},
			},
		},
		{
			"status without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateStatus{Active},
			Action{Status: Active},
			Action{Status: Active},
			false,
			service.NewMutations("", ""),
		},
		{
			"project with change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateProject{"bar"},
			Action{Project: "foo"},
			Action{Project: "bar"},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From:        map[string]any{"project": "foo"},
				To:          map[string]any{"project": "bar"},
			},
		},
		{
			"project without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateProject{"foo"},
			Action{Project: "foo"},
			Action{Project: "foo"},
			false,
			service.NewMutations("", ""),
		},
		{
			"add tag with change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateTagsAdd{"bar"},
			Action{Tags: []string{"foo"}},
			Action{Tags: []string{"foo", "bar"}},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From:        map[string]any{"tags": []string{"foo"}},
				To:          map[string]any{"tags": []string{"foo", "bar"}},
			},
		},
		{
			"add tag without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateTagsAdd{"foo"},
			Action{Tags: []string{"foo"}},
			Action{Tags: []string{"foo"}},
			false,
			service.NewMutations("", ""),
		},
		{
			"remove tag with change",
			func() *service.Mutations {
				return &service.Mutations{
					Transaction: "",
					Id:          "",
					From: map[string]any{
						"tags": []string{"foo", "bar", "fiz"},
					},
					To: map[string]any{"tags": []string{"foo", "bar"}},
				}
			},
			&MutateTagsRemove{"foo"},
			Action{Tags: []string{"foo", "bar"}},
			Action{Tags: []string{"bar"}},
			true,
			&service.Mutations{
				Transaction: "",
				Id:          "",
				From: map[string]any{
					"tags": []string{"foo", "bar", "fiz"},
				},
				To: map[string]any{"tags": []string{"bar"}},
			},
		},
		{
			"remove tag without change",
			func() *service.Mutations { return service.NewMutations("", "") },
			&MutateTagsRemove{"fiz"},
			Action{Tags: []string{"foo", "bar"}},
			Action{Tags: []string{"foo", "bar"}},
			false,
			service.NewMutations("", ""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := tc.mutations()
			got := tc.mutator.Mutate(&tc.in, now, ms)
			if !reflect.DeepEqual(&tc.after, &tc.in) {
				t.Errorf("expected mutation to \n%+v but got \n%+v", tc.after, tc.in)
			}
			if got != tc.expect {
				t.Errorf("expected %t but got %t", tc.expect, got)
			}
			if !reflect.DeepEqual(ms, tc.expectedMutations) {
				t.Errorf(
					"expected mutations \n%+v, but got \n%+v",
					tc.expectedMutations,
					ms,
				)
			}
		})
	}
}
