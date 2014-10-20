package types

type Page interface {
	Destroy() error
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookie(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(width, height int) error
	Screenshot(filename string) error
	Title() (string, error)
	HTML() (string, error)
	RunScript(body string, arguments map[string]interface{}, result interface{}) error
	Forward() error
	Back() error
	Refresh() error
	Find(selector string) Selection
	FindByXPath(selector string) Selection
	FindByLink(text string) Selection
	FindByLabel(text string) Selection
	All(selector string) MultiSelection
	AllByXPath(selector string) MultiSelection
	AllByLink(text string) MultiSelection
	AllByLabel(text string) MultiSelection
}
