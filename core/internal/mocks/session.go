package mocks

import (
	"encoding/json"

	"github.com/sclevine/agouti/api"
)

type Session struct {
	GetElementCall struct {
		Selector      api.Selector
		ReturnElement *api.Element
		Err           error
	}

	GetElementsCall struct {
		Selector       api.Selector
		ReturnElements []*api.Element
		Err            error
	}

	GetActiveElementCall struct {
		ReturnElement *api.Element
		Err           error
	}

	DeleteCall struct {
		Called bool
		Err    error
	}

	GetWindowCall struct {
		ReturnWindow *api.Window
		Err          error
	}

	GetWindowsCall struct {
		ReturnWindows []*api.Window
		Err           error
	}

	SetWindowCall struct {
		Window *api.Window
		Err    error
	}

	SetWindowByNameCall struct {
		Name string
		Err  error
	}

	DeleteWindowCall struct {
		Called bool
		Err    error
	}

	GetScreenshotCall struct {
		ReturnImage []byte
		Err         error
	}

	SetCookieCall struct {
		Cookie map[string]interface{}
		Err    error
	}

	DeleteCookieCall struct {
		Name string
		Err  error
	}

	DeleteCookiesCall struct {
		Called bool
		Err    error
	}

	GetURLCall struct {
		ReturnURL string
		Err       error
	}

	SetURLCall struct {
		URL string
		Err error
	}

	GetTitleCall struct {
		ReturnTitle string
		Err         error
	}

	GetSourceCall struct {
		ReturnSource string
		Err          error
	}

	DoubleClickCall struct {
		Called bool
		Err    error
	}

	MoveToCall struct {
		Element *api.Element
		Point   api.Offset
		Err     error
	}

	FrameCall struct {
		Frame *api.Element
		Err   error
	}

	FrameParentCall struct {
		Called bool
		Err    error
	}

	ExecuteCall struct {
		Body      string
		Arguments []interface{}
		Result    string
		Err       error
	}

	ForwardCall struct {
		Called bool
		Err    error
	}

	BackCall struct {
		Called bool
		Err    error
	}

	RefreshCall struct {
		Called bool
		Err    error
	}

	GetAlertTextCall struct {
		ReturnText string
		Err        error
	}

	SetAlertTextCall struct {
		Text string
		Err  error
	}

	AcceptAlertCall struct {
		Called bool
		Err    error
	}

	DismissAlertCall struct {
		Called bool
		Err    error
	}

	NewLogsCall struct {
		LogType    string
		ReturnLogs []api.Log
		Err        error
	}

	GetLogTypesCall struct {
		ReturnTypes []string
		Err         error
	}
}

func (s *Session) Delete() error {
	s.DeleteCall.Called = true
	return s.DeleteCall.Err
}

func (s *Session) GetElement(selector api.Selector) (*api.Element, error) {
	s.GetElementCall.Selector = selector
	return s.GetElementCall.ReturnElement, s.GetElementCall.Err
}

func (s *Session) GetElements(selector api.Selector) ([]*api.Element, error) {
	s.GetElementsCall.Selector = selector
	return s.GetElementsCall.ReturnElements, s.GetElementsCall.Err
}

func (s *Session) GetActiveElement() (*api.Element, error) {
	return s.GetActiveElementCall.ReturnElement, s.GetActiveElementCall.Err
}

func (s *Session) GetWindow() (*api.Window, error) {
	return s.GetWindowCall.ReturnWindow, s.GetWindowCall.Err
}

func (s *Session) GetWindows() ([]*api.Window, error) {
	return s.GetWindowsCall.ReturnWindows, s.GetWindowsCall.Err
}

func (s *Session) SetWindow(window *api.Window) error {
	s.SetWindowCall.Window = window
	return s.SetWindowCall.Err
}

func (s *Session) SetWindowByName(name string) error {
	s.SetWindowByNameCall.Name = name
	return s.SetWindowByNameCall.Err
}

func (s *Session) DeleteWindow() error {
	s.DeleteWindowCall.Called = true
	return s.DeleteWindowCall.Err
}

func (s *Session) GetScreenshot() ([]byte, error) {
	return s.GetScreenshotCall.ReturnImage, s.GetScreenshotCall.Err
}

func (s *Session) SetCookie(cookie map[string]interface{}) error {
	s.SetCookieCall.Cookie = cookie
	return s.SetCookieCall.Err
}

func (s *Session) DeleteCookie(name string) error {
	s.DeleteCookieCall.Name = name
	return s.DeleteCookieCall.Err
}

func (s *Session) DeleteCookies() error {
	s.DeleteCookiesCall.Called = true
	return s.DeleteCookiesCall.Err
}

func (s *Session) GetURL() (string, error) {
	return s.GetURLCall.ReturnURL, s.GetURLCall.Err
}

func (s *Session) SetURL(url string) error {
	s.SetURLCall.URL = url
	return s.SetURLCall.Err
}

func (s *Session) GetTitle() (string, error) {
	return s.GetTitleCall.ReturnTitle, s.GetTitleCall.Err
}

func (s *Session) GetSource() (string, error) {
	return s.GetSourceCall.ReturnSource, s.GetSourceCall.Err
}

func (s *Session) DoubleClick() error {
	s.DoubleClickCall.Called = true
	return s.DoubleClickCall.Err
}

func (s *Session) MoveTo(element *api.Element, point api.Offset) error {
	s.MoveToCall.Element = element
	s.MoveToCall.Point = point
	return s.MoveToCall.Err
}

func (s *Session) Frame(frame *api.Element) error {
	s.FrameCall.Frame = frame
	return s.FrameCall.Err
}

func (s *Session) FrameParent() error {
	s.FrameParentCall.Called = true
	return s.FrameParentCall.Err
}

func (s *Session) Execute(body string, arguments []interface{}, result interface{}) error {
	s.ExecuteCall.Body = body
	s.ExecuteCall.Arguments = arguments
	json.Unmarshal([]byte(s.ExecuteCall.Result), result)
	return s.ExecuteCall.Err
}

func (s *Session) Forward() error {
	s.ForwardCall.Called = true
	return s.ForwardCall.Err
}

func (s *Session) Back() error {
	s.BackCall.Called = true
	return s.BackCall.Err
}

func (s *Session) Refresh() error {
	s.RefreshCall.Called = true
	return s.RefreshCall.Err
}

func (s *Session) GetAlertText() (string, error) {
	return s.GetAlertTextCall.ReturnText, s.GetAlertTextCall.Err
}

func (s *Session) SetAlertText(text string) error {
	s.SetAlertTextCall.Text = text
	return s.SetAlertTextCall.Err
}

func (s *Session) AcceptAlert() error {
	s.AcceptAlertCall.Called = true
	return s.AcceptAlertCall.Err
}

func (s *Session) DismissAlert() error {
	s.DismissAlertCall.Called = true
	return s.DismissAlertCall.Err
}

func (s *Session) NewLogs(logType string) ([]api.Log, error) {
	s.NewLogsCall.LogType = logType
	return s.NewLogsCall.ReturnLogs, s.NewLogsCall.Err
}

func (s *Session) GetLogTypes() ([]string, error) {
	return s.GetLogTypesCall.ReturnTypes, s.GetLogTypesCall.Err
}
