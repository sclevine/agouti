package mocks

import "github.com/sclevine/agouti"

type Page struct {
	DestroyCall struct {
		Called bool
		Err    error
	}

	NavigateCall struct {
		URL string
		Err error
	}

	SetCookieCall struct {
		Cookie agouti.Cookie
		Err    error
	}

	DeleteCookieCall struct {
		Name string
		Err  error
	}

	ClearCookiesCall struct {
		Called bool
		Err    error
	}

	SizeCall struct {
		Width  int
		Height int
		Err    error
	}

	ScreenshotCall struct {
		Filename string
		Err      error
	}

	RunScriptCall struct {
		Body      string
		Arguments map[string]interface{}
		Result    interface{}
		Err       error
	}

	EnterPopupTextCall struct {
		Text string
		Err  error
	}

	ConfirmPopupCall struct {
		Called bool
		Err    error
	}

	CancelPopupCall struct {
		Called bool
		Err    error
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

	SwitchToParentFrameCall struct {
		Called bool
		Err    error
	}

	SwitchToRootFrameCall struct {
		Called bool
		Err    error
	}

	SwitchToWindowCall struct {
		Name string
		Err  error
	}

	NextWindowCall struct {
		Called bool
		Err    error
	}

	CloseWindowCall struct {
		Called bool
		Err    error
	}
}

func (p *Page) Destroy() error {
	p.DestroyCall.Called = true
	return p.DestroyCall.Err
}

func (p *Page) Navigate(url string) error {
	p.NavigateCall.URL = url
	return p.NavigateCall.Err
}

func (p *Page) SetCookie(cookie agouti.Cookie) error {
	p.SetCookieCall.Cookie = cookie
	return p.SetCookieCall.Err
}

func (p *Page) DeleteCookie(name string) error {
	p.DeleteCookieCall.Name = name
	return p.DeleteCookieCall.Err
}

func (p *Page) ClearCookies() error {
	p.ClearCookiesCall.Called = true
	return p.ClearCookiesCall.Err
}

func (p *Page) Size(width, height int) error {
	p.SizeCall.Width = width
	p.SizeCall.Height = height
	return p.SizeCall.Err
}

func (p *Page) Screenshot(filename string) error {
	p.ScreenshotCall.Filename = filename
	return p.ScreenshotCall.Err
}

func (p *Page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
	p.RunScriptCall.Body = body
	p.RunScriptCall.Arguments = arguments
	p.RunScriptCall.Result = result
	return p.RunScriptCall.Err
}

func (p *Page) EnterPopupText(text string) error {
	p.EnterPopupTextCall.Text = text
	return p.EnterPopupTextCall.Err
}

func (p *Page) ConfirmPopup() error {
	p.ConfirmPopupCall.Called = true
	return p.ConfirmPopupCall.Err
}

func (p *Page) CancelPopup() error {
	p.CancelPopupCall.Called = true
	return p.CancelPopupCall.Err
}

func (p *Page) Back() error {
	p.BackCall.Called = true
	return p.BackCall.Err
}

func (p *Page) Forward() error {
	p.ForwardCall.Called = true
	return p.ForwardCall.Err
}

func (p *Page) Refresh() error {
	p.RefreshCall.Called = true
	return p.RefreshCall.Err
}

func (p *Page) SwitchToParentFrame() error {
	p.SwitchToParentFrameCall.Called = true
	return p.SwitchToParentFrameCall.Err
}

func (p *Page) SwitchToRootFrame() error {
	p.SwitchToRootFrameCall.Called = true
	return p.SwitchToRootFrameCall.Err
}

func (p *Page) SwitchToWindow(name string) error {
	p.SwitchToWindowCall.Name = name
	return p.SwitchToWindowCall.Err
}

func (p *Page) NextWindow() error {
	p.NextWindowCall.Called = true
	return p.NextWindowCall.Err
}

func (p *Page) CloseWindow() error {
	p.CloseWindowCall.Called = true
	return p.CloseWindowCall.Err
}
