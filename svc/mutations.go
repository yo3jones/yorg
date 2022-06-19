package svc

type Mutations struct {
	Transaction string
	Id          string
	From        map[string]any
	To          map[string]any
}

func NewMutations(transaction, id string) *Mutations {
	return &Mutations{
		Transaction: transaction,
		Id:          id,
		From:        make(map[string]any),
		To:          make(map[string]any),
	}
}

func (m *Mutations) Add(field string, from, to any) {
	_, exists := m.From[field]

	if !exists {
		m.From[field] = from
	}
	m.To[field] = to
}
