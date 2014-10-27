package selection

type MultiSelection struct {
	*Selection
}

func (m *MultiSelection) At(index int) *Selection {
	return m.modifiedSelection(false, true, index)
}

func (m *MultiSelection) Single() *Selection {
	return m.modifiedSelection(true, false)
}

func (m *MultiSelection) modifiedSelection(single, indexed bool, index ...int) *Selection {
	last := len(m.selectors) - 1

	if last < 0 {
		return &Selection{m.Client, nil}
	}

	newSelector := m.selectors[last]
	newSelector.Single = single
	newSelector.Indexed = indexed
	if len(index) > 0 {
		newSelector.Index = index[0]
	}
	return &Selection{m.Client, appendSelector(m.selectors[:last], newSelector)}
}
