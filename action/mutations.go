package action

import (
	"strings"
	"time"

	"github.com/yo3jones/yorg/service"
)

type MutateType struct {
	value ActionType
}

func (m *MutateType) Mutate(a *Action, _ time.Time, ms *service.Mutations) bool {
	if a.Type == m.value {
		return false
	}

	ms.Add("type", a.Type, m.value)

	a.Type = m.value

	return true
}

type MutateAction struct {
	value string
}

func (m *MutateAction) Mutate(a *Action, _ time.Time, ms *service.Mutations) bool {
	if a.Action == m.value {
		return false
	}

	ms.Add("action", a.Action, m.value)

	a.Action = m.value

	return true
}

type MutateStatus struct {
	value ActionStatus
}

func (m *MutateStatus) Mutate(a *Action, t time.Time, ms *service.Mutations) bool {
	if a.Status == m.value {
		return false
	}

	ms.Add("status", a.Status, m.value)

	if m.value == Complete {
		ms.Add("completedAt", a.CompletedAt, t)
		a.CompletedAt = t
	} else {
		ms.Add("completedAt", a.CompletedAt, time.Time{})
		a.CompletedAt = time.Time{}
	}

	a.Status = m.value

	return true
}

type MutateProject struct {
	value string
}

func (m *MutateProject) Mutate(a *Action, _ time.Time, ms *service.Mutations) bool {
	if a.Project == m.value {
		return false
	}

	ms.Add("project", a.Project, m.value)

	a.Project = m.value

	return true
}

type MutateTagsAdd struct {
	value string
}

func (m *MutateTagsAdd) Mutate(a *Action, _ time.Time, ms *service.Mutations) bool {
	for _, tag := range a.Tags {
		if strings.ToLower(tag) == strings.ToLower(m.value) {
			return false
		}
	}

	updatedTags := append([]string{}, a.Tags...)
	updatedTags = append(updatedTags, strings.ToLower(m.value))
	ms.Add("tags", a.Tags, updatedTags)

	a.Tags = updatedTags

	return true
}

type MutateTagsRemove struct {
	value string
}

func (m *MutateTagsRemove) Mutate(a *Action, _ time.Time, ms *service.Mutations) bool {
	index := -1
	for i, tag := range a.Tags {
		if strings.ToLower(tag) == strings.ToLower(m.value) {
			index = i
			break
		}
	}

	if index < 0 {
		return false
	}

	updatedTags := append([]string{}, a.Tags[:index]...)
	updatedTags = append(updatedTags, a.Tags[index+1:]...)
	ms.Add("tags", a.Tags, updatedTags)

	a.Tags = updatedTags

	return true
}
