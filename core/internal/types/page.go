package types

type Page interface {
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookie(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(width, height int) error
	Screenshot(filename string) error
	Title() (string, error)
	RunScript(body string, arguments map[string]interface{}, result interface{}) error
	Forward() error
	Back() error
	Refresh() error
	Find(selector string) Selection
	FindXPath(selector string) Selection
	FindByLabel(text string) Selection
}
