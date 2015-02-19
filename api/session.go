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
	Bus
}

type Bus interface {
	Send(endpoint, method string, body, result interface{}) error
}

func Open(url string, capabilities map[string]interface{}) (*Session, error) {
	busClient, err := bus.Connect(url, capabilities)
	if err != nil {
		return nil, err
	}
	return &Session{busClient}, nil
}

func (s *Session) Delete() error {
	return s.Send("", "DELETE", nil, nil)
}

func (s *Session) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := s.Send("element", "POST", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := s.Send("elements", "POST", selector, &results); err != nil {
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

	if err := s.Send("element/active", "POST", nil, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetWindow() (*Window, error) {
	var windowID string
	if err := s.Send("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{windowID, s}, nil
}

func (s *Session) GetWindows() ([]*Window, error) {
	var windowsID []string
	if err := s.Send("window_handles", "GET", nil, &windowsID); err != nil {
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

	return s.Send("window", "POST", request, nil)
}

func (s *Session) SetWindowByName(name string) error {
	request := struct {
		Name string `json:"name"`
	}{name}

	return s.Send("window", "POST", request, nil)
}

func (s *Session) DeleteWindow() error {
	if err := s.Send("window", "DELETE", nil, nil); err != nil {
		return err
	}
	return nil
}

func (s *Session) SetCookie(cookie map[string]interface{}) error {
	request := struct {
		Cookie map[string]interface{} `json:"cookie"`
	}{cookie}

	return s.Send("cookie", "POST", request, nil)
}

func (s *Session) DeleteCookie(cookieName string) error {
	return s.Send("cookie/"+cookieName, "DELETE", nil, nil)
}

func (s *Session) DeleteCookies() error {
	return s.Send("cookie", "DELETE", nil, nil)
}

func (s *Session) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := s.Send("screenshot", "GET", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (s *Session) GetURL() (string, error) {
	var url string
	if err := s.Send("url", "GET", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (s *Session) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return s.Send("url", "POST", request, nil)
}

func (s *Session) GetTitle() (string, error) {
	var title string
	if err := s.Send("title", "GET", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (s *Session) GetSource() (string, error) {
	var source string
	if err := s.Send("source", "GET", nil, &source); err != nil {
		return "", err
	}

	return source, nil
}

func (s *Session) DoubleClick() error {
	return s.Send("doubleclick", "POST", nil, nil)
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

	return s.Send("moveto", "POST", request, nil)
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

	return s.Send("frame", "POST", request, nil)
}

func (s *Session) FrameParent() error {
	return s.Send("frame/parent", "POST", nil, nil)
}

func (s *Session) Execute(body string, arguments []interface{}, result interface{}) error {
	if arguments == nil {
		arguments = []interface{}{}
	}

	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := s.Send("execute", "POST", request, result); err != nil {
		return err
	}

	return nil
}

func (s *Session) Forward() error {
	return s.Send("forward", "POST", nil, nil)
}

func (s *Session) Back() error {
	return s.Send("back", "POST", nil, nil)
}

func (s *Session) Refresh() error {
	return s.Send("refresh", "POST", nil, nil)
}

func (s *Session) GetAlertText() (string, error) {
	var text string
	if err := s.Send("alert_text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Session) SetAlertText(text string) error {
	request := struct {
		Text string `json:"text"`
	}{text}
	return s.Send("alert_text", "POST", request, nil)
}

func (s *Session) AcceptAlert() error {
	return s.Send("accept_alert", "POST", nil, nil)
}

func (s *Session) DismissAlert() error {
	return s.Send("dismiss_alert", "POST", nil, nil)
}

func (s *Session) NewLogs(logType string) ([]Log, error) {
	request := struct {
		Type string `json:"type"`
	}{logType}

	var logs []Log
	if err := s.Send("log", "POST", request, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Session) GetLogTypes() ([]string, error) {
	var types []string
	if err := s.Send("log/types", "GET", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}

func (s *Session) sendElement(id, endpoint, method string, body, result interface{}) error {
	return s.Send(path.Join("element", id, endpoint), method, body, result)
}

func (s *Session) sendWindow(id, endpoint, method string, body interface{}, result interface{}) error {
	return s.Send(path.Join("window", id, endpoint), method, body, result)
}
