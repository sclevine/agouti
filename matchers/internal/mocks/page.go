package mocks

import "github.com/sclevine/agouti"

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

	ReadAllLogsCall struct {
		LogType    string
		ReturnLogs []agouti.Log
		Err        error
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

func (p *Page) ReadAllLogs(logType string) ([]agouti.Log, error) {
	p.ReadAllLogsCall.LogType = logType
	return p.ReadAllLogsCall.ReturnLogs, p.ReadAllLogsCall.Err
}
