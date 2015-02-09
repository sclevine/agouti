package mocks

type Selection struct {
	SwitchToFrameCall struct {
		Called bool
		Err    error
	}

	ClickCall struct {
		Called bool
		Err    error
	}

	DoubleClickCall struct {
		Called bool
		Err    error
	}

	FillCall struct {
		Text string
		Err  error
	}

	CheckCall struct {
		Called bool
		Err    error
	}

	UncheckCall struct {
		Called bool
		Err    error
	}

	SelectCall struct {
		Text string
		Err  error
	}

	SubmitCall struct {
		Called bool
		Err    error
	}
}

func (s *Selection) SwitchToFrame() error {
	s.SwitchToFrameCall.Called = true
	return s.SwitchToFrameCall.Err
}

func (s *Selection) Click() error {
	s.ClickCall.Called = true
	return s.ClickCall.Err
}

func (s *Selection) DoubleClick() error {
	s.DoubleClickCall.Called = true
	return s.DoubleClickCall.Err
}

func (s *Selection) Fill(text string) error {
	s.FillCall.Text = text
	return s.FillCall.Err
}

func (s *Selection) Check() error {
	s.CheckCall.Called = true
	return s.CheckCall.Err
}

func (s *Selection) Uncheck() error {
	s.UncheckCall.Called = true
	return s.UncheckCall.Err
}

func (s *Selection) Select(text string) error {
	s.SelectCall.Text = text
	return s.SelectCall.Err
}

func (s *Selection) Submit() error {
	s.SubmitCall.Called = true
	return s.SubmitCall.Err
}
