package mocks

import "github.com/sclevine/agouti/core"

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

	ReadLogsCall struct {
		All        bool
		LogType    string
		ReturnLogs []core.Log
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

func (p *Page) ReadLogs(logType string, all ...bool) ([]core.Log, error) {
	p.ReadLogsCall.LogType = logType
	p.ReadLogsCall.All = len(all) > 0 && all[0]
	return p.ReadLogsCall.ReturnLogs, p.ReadLogsCall.Err
}
