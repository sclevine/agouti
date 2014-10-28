package selection

type MultiSelection struct {
	*Selection
}

func (m *MultiSelection) At(index int) *Selection {
	last := len(m.selectors) - 1

	if last < 0 {
		return &Selection{m.Client, nil}
	}

	newSelector := m.selectors[last]
	newSelector.Indexed = true
	newSelector.Index = index

	return &Selection{m.Client, appendSelector(m.selectors[:last], newSelector)}
}

func (m *MultiSelection) Single() *Selection {
	last := len(m.selectors) - 1

	if last < 0 {
		return &Selection{m.Client, nil}
	}

	newSelector := m.selectors[last]
	newSelector.Single = true

	return &Selection{m.Client, appendSelector(m.selectors[:last], newSelector)}
}
