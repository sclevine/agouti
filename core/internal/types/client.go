package types

type Client interface {
	DeleteSession() error
	GetWindow() (Window, error)
	GetScreenshot() ([]byte, error)
	SetCookie(cookie interface{}) error
	DeleteCookie(name string) error
	DeleteCookies() error
	GetURL() (string, error)
	SetURL(url string) error
	GetTitle() (string, error)
	GetSource() (string, error)
	GetElement(selector Selector) (Element, error)
	GetActiveElement() (Element, error)
	GetElements(selector Selector) ([]Element, error)
	DoubleClick() error
	MoveTo(element Element, point Point) error
	Frame(frame Element) error
	FrameParent() error
	Execute(body string, arguments []interface{}, result interface{}) error
	Forward() error
	Back() error
	Refresh() error
	GetAlertText() (string, error)
	SetAlertText(text string) error
	NewLogs(logType string) ([]Log, error)
	GetLogTypes() ([]string, error)
}
