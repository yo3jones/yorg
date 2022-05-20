package action

import (
	"strings"
)

type IdEquals struct {
	value string
}

func (m *IdEquals) Match(a *Action) bool {
	return a.Id == m.value
}

type TypeEquals struct {
	value ActionType
}

func (m *TypeEquals) Match(a *Action) bool {
	return a.Type == m.value
}

type StatusEquals struct {
	value ActionStatus
}

func (m *StatusEquals) Match(a *Action) bool {
	return a.Status == m.value
}

type ProjectEquals struct {
	value string
}

func (m *ProjectEquals) Match(a *Action) bool {
	return a.Project == m.value
}

type HasTag struct {
	value string
}

func (m *HasTag) Match(a *Action) bool {
	for _, tag := range a.Tags {
		if strings.ToLower(tag) == strings.ToLower(m.value) {
			return true
		}
	}
	return false
}

type DoesNotHaveTag struct {
	value string
}

func (m *DoesNotHaveTag) Match(a *Action) bool {
	for _, tag := range a.Tags {
		if strings.ToLower(tag) == strings.ToLower(m.value) {
			return false
		}
	}
	return true
}
