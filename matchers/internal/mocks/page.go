package mocks

type Page struct {
	TitleCall struct {
		ReturnTitle string
		Err         error
	}

	PopupTextCall struct {
		ReturnText string
		Err        error
	}

	URLCall struct {
		ReturnURL string
		Err       error
	}
}

func (p *Page) Title() (string, error) {
	return p.TitleCall.ReturnTitle, p.TitleCall.Err
}

func (p *Page) PopupText() (string, error) {
	return p.PopupTextCall.ReturnText, p.PopupTextCall.Err
}

func (p *Page) URL() (string, error) {
	return p.URLCall.ReturnURL, p.URLCall.Err
}
