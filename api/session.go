package api

import (
	"encoding/base64"
	"path"

	"github.com/sclevine/agouti/api/internal/bus"
)

type Log struct {
	Message   string
	Level     string
	Timestamp int64
}

type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

type Session struct {
	Bus busSender
}

type busSender interface {
	Send(endpoint, method string, body interface{}, result ...interface{}) error
}

func Open(url string, capabilities map[string]interface{}) (*Session, error) {
	busClient, err := bus.Connect(url, capabilities)
	if err != nil {
		return nil, err
	}
	return &Session{busClient}, nil
}

func (s *Session) Delete() error {
	return s.Bus.Send("", "DELETE", nil)
}

func (s *Session) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := s.Bus.Send("element", "POST", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := s.Bus.Send("elements", "POST", selector, &results); err != nil {
		return nil, err
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, s})
	}

	return elements, nil
}

func (s *Session) GetActiveElement() (*Element, error) {
	var result struct{ Element string }

	if err := s.Bus.Send("element/active", "POST", nil, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetWindow() (*Window, error) {
	var windowID string
	if err := s.Bus.Send("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{windowID, s}, nil
}

func (s *Session) GetWindows() ([]*Window, error) {
	var windowsID []string
	if err := s.Bus.Send("window_handles", "GET", nil, &windowsID); err != nil {
		return nil, err
	}

	var windows []*Window
	for _, windowID := range windowsID {
		windows = append(windows, &Window{windowID, s})
	}
	return windows, nil
}

func (s *Session) SetWindow(window *Window) error {
	request := struct {
		Name string `json:"name"`
	}{window.ID}

	return s.Bus.Send("window", "POST", request)
}

func (s *Session) SetWindowByName(name string) error {
	request := struct {
		Name string `json:"name"`
	}{name}

	return s.Bus.Send("window", "POST", request)
}

func (s *Session) DeleteWindow() error {
	if err := s.Bus.Send("window", "DELETE", nil, nil); err != nil {
		return err
	}
	return nil
}

func (s *Session) SetCookie(cookie map[string]interface{}) error {
	request := struct {
		Cookie map[string]interface{} `json:"cookie"`
	}{cookie}

	return s.Bus.Send("cookie", "POST", request)
}

func (s *Session) DeleteCookie(cookieName string) error {
	return s.Bus.Send("cookie/"+cookieName, "DELETE", nil)
}

func (s *Session) DeleteCookies() error {
	return s.Bus.Send("cookie", "DELETE", nil)
}

func (s *Session) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := s.Bus.Send("screenshot", "GET", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (s *Session) GetURL() (string, error) {
	var url string
	if err := s.Bus.Send("url", "GET", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (s *Session) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return s.Bus.Send("url", "POST", request)
}

func (s *Session) GetTitle() (string, error) {
	var title string
	if err := s.Bus.Send("title", "GET", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (s *Session) GetSource() (string, error) {
	var source string
	if err := s.Bus.Send("source", "GET", nil, &source); err != nil {
		return "", err
	}

	return source, nil
}

func (s *Session) DoubleClick() error {
	return s.Bus.Send("doubleclick", "POST", nil)
}

func (s *Session) MoveTo(region *Element, offset Offset) error {
	request := map[string]interface{}{}

	if region != nil {
		// TODO: return error if not an element
		request["element"] = region.ID
	}

	if offset != nil {
		if xoffset, present := offset.x(); present {
			request["xoffset"] = xoffset
		}

		if yoffset, present := offset.y(); present {
			request["yoffset"] = yoffset
		}
	}

	return s.Bus.Send("moveto", "POST", request)
}

func (s *Session) Frame(frame *Element) error {
	var elementID interface{}

	if frame != nil {
		elementID = struct {
			Element string `json:"ELEMENT"`
		}{frame.ID}
	}

	request := struct {
		ID interface{} `json:"id"`
	}{elementID}

	return s.Bus.Send("frame", "POST", request)
}

func (s *Session) FrameParent() error {
	return s.Bus.Send("frame/parent", "POST", nil)
}

func (s *Session) Execute(body string, arguments []interface{}, result interface{}) error {
	if arguments == nil {
		arguments = []interface{}{}
	}

	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := s.Bus.Send("execute", "POST", request, result); err != nil {
		return err
	}

	return nil
}

func (s *Session) Forward() error {
	return s.Bus.Send("forward", "POST", nil)
}

func (s *Session) Back() error {
	return s.Bus.Send("back", "POST", nil)
}

func (s *Session) Refresh() error {
	return s.Bus.Send("refresh", "POST", nil)
}

func (s *Session) GetAlertText() (string, error) {
	var text string
	if err := s.Bus.Send("alert_text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Session) SetAlertText(text string) error {
	request := struct {
		Text string `json:"text"`
	}{text}
	return s.Bus.Send("alert_text", "POST", request)
}

func (s *Session) AcceptAlert() error {
	return s.Bus.Send("accept_alert", "POST", nil)
}

func (s *Session) DismissAlert() error {
	return s.Bus.Send("dismiss_alert", "POST", nil)
}

func (s *Session) NewLogs(logType string) ([]Log, error) {
	request := struct {
		Type string `json:"type"`
	}{logType}

	var logs []Log
	if err := s.Bus.Send("log", "POST", request, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Session) GetLogTypes() ([]string, error) {
	var types []string
	if err := s.Bus.Send("log/types", "GET", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}

func (s *Session) sendElement(id, endpoint, method string, body interface{}, result ...interface{}) error {
	return s.Bus.Send(path.Join("element", id, endpoint), method, body, result...)
}

func (s *Session) sendWindow(id, endpoint, method string, body interface{}, result ...interface{}) error {
	return s.Bus.Send(path.Join("window", id, endpoint), method, body, result...)
}
