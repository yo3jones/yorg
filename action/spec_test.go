package action

import (
	"fmt"
	"reflect"
	"testing"
)

func TestActionType(t *testing.T) {
	type test struct {
		name       string
		actionType ActionType
		str        string
	}

	tests := []test{
		{"unknown", UnknownActionType, "unknown"},
		{"inbox", Inbox, "inbox"},
		{"next-action", NextAction, "next-action"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStr := tc.actionType.String()
			if gotStr != tc.str {
				t.Errorf("excpected %s but got %s", tc.str, gotStr)
			}
			gotType := ActionTypeFromString(tc.str)
			if gotType != tc.actionType {
				t.Errorf("expected %s but got %s", tc.actionType, gotType)
			}
			var gotUnmarshal ActionType
			err1 := gotUnmarshal.UnmarshalJSON(
				[]byte(fmt.Sprintf(`"%s"`, tc.str)),
			)
			if err1 != nil {
				t.Error(err1)
			}
			if gotUnmarshal != tc.actionType {
				t.Errorf("expected %s but got %s", tc.actionType, gotUnmarshal)
			}
			gotMarshal, err2 := tc.actionType.MarshalJSON()
			if err2 != nil {
				t.Error(err2)
			}
			if reflect.DeepEqual(gotMarshal, []byte(tc.str)) {
				t.Errorf("expected %s but ot %s", tc.str, gotMarshal)
			}
		})
	}

	t.Run("with bad input", func(t *testing.T) {
		var got ActionType
		err := got.UnmarshalJSON([]byte("bad"))
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})
}

func TestActionStatus(t *testing.T) {
	type test struct {
		name         string
		actionStatus ActionStatus
		str          string
	}

	tests := []test{
		{"unknown", UnknownActionStatus, "unknown"},
		{"active", Active, "active"},
		{"complete", Complete, "complete"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStr := tc.actionStatus.String()
			if gotStr != tc.str {
				t.Errorf("excpected %s but got %s", tc.str, gotStr)
			}
			gotStatus := ActionStatusFromString(tc.str)
			if gotStatus != tc.actionStatus {
				t.Errorf("expected %s but got %s", tc.actionStatus, gotStatus)
			}
			var gotUnmarshal ActionStatus
			err1 := gotUnmarshal.UnmarshalJSON(
				[]byte(fmt.Sprintf(`"%s"`, tc.str)),
			)
			if err1 != nil {
				t.Error(err1)
			}
			if gotUnmarshal != tc.actionStatus {
				t.Errorf("expected %s but got %s", tc.actionStatus, gotUnmarshal)
			}
			gotMarshal, err2 := tc.actionStatus.MarshalJSON()
			if err2 != nil {
				t.Error(err2)
			}
			if reflect.DeepEqual(gotMarshal, []byte(tc.str)) {
				t.Errorf("expected %s but ot %s", tc.str, gotMarshal)
			}
		})
	}

	t.Run("with bad input", func(t *testing.T) {
		var got ActionStatus
		err := got.UnmarshalJSON([]byte("bad"))
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})
}
