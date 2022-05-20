package action

import (
	"time"
)

func stringLesser[T string](i, j T, desc bool) int {
	if i == j {
		return 0
	}
	var result int
	if i < j {
		result = -1
	} else {
		result = 1
	}
	if desc {
		return result * -1
	}
	return result
}

func timeLesser(i, j time.Time, desc bool) int {
	if i == j {
		return 0
	}
	var result int
	if i.Before(j) {
		result = -1
	} else {
		result = 1
	}
	if desc {
		return result * -1
	}

	return result
}

type OrderById struct {
	desc bool
}

func (l *OrderById) Less(i, j *Action) int {
	return stringLesser(i.Id, j.Id, l.desc)
}

type OrderByCreatedAt struct {
	desc bool
}

func (l *OrderByCreatedAt) Less(i, j *Action) int {
	return timeLesser(i.CreatedAt, j.CreatedAt, l.desc)
}
