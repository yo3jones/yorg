package action

import (
	"encoding/json"
	"time"
)

type Action struct {
	Id          string       `json:"id"`
	Type        ActionType   `json:"type"`
	Action      string       `json:"action"`
	Status      ActionStatus `json:"status"`
	Project     string       `json:"project"`
	Tags        []string     `json:"tags"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	CompletedAt time.Time    `json:"completedAt,omitempty"`
}

type ActionType int

const (
	UnknownActionType ActionType = iota
	Inbox
	NextAction
)

func (at ActionType) String() string {
	switch at {
	case Inbox:
		return "inbox"
	case NextAction:
		return "next-action"
	}
	return "unknown"
}

func ActionTypeFromString(s string) ActionType {
	switch s {
	case "inbox":
		return Inbox
	case "next-action":
		return NextAction
	}

	return UnknownActionType
}

func (at *ActionType) UnmarshalJSON(b []byte) error {
	var s string

	err1 := json.Unmarshal(b, &s)
	if err1 != nil {
		return err1
	}

    *at = ActionTypeFromString(s)

	return nil
}

func (at ActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.String())
}

type ActionStatus int

const (
	UnknownActionStatus ActionStatus = iota
	Active
	Complete
)

func (as ActionStatus) String() string {
	switch as {
	case Active:
		return "active"
	case Complete:
		return "complete"
	}
	return "unknown"
}

func ActionStatusFromString(s string) ActionStatus {
	switch s {
	case "active":
		return Active
	case "complete":
		return Complete
	}

	return UnknownActionStatus
}

func (as ActionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(as.String())
}

func (as *ActionStatus) UnmarshalJSON(b []byte) error {
	var s string

	err1 := json.Unmarshal(b, &s)
	if err1 != nil {
		return err1
	}

    *as = ActionStatusFromString(s)

	return nil
}
