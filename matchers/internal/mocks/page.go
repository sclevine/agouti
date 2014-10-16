package mocks

type Page struct {
	TitleCall struct {
		ReturnTitle string
		Err         error
	}
}

func (p *Page) Title() (string, error) {
	return p.TitleCall.ReturnTitle, p.TitleCall.Err
}
