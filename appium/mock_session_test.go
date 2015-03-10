package appium

import "github.com/sclevine/agouti/api"

type mockMobileSession struct {
	PerformTouchCall struct {
		Selector      api.Selector
		ReturnElement *api.Element
		Err           error
	}

	LaunchAppCall struct {
		Err error
	}

	CloseAppCall struct {
		Err error
	}

	InstallAppCall struct {
		Err error
	}
}

func (ms *mockMobileSession) LaunchApp() error {
	return ms.LaunchAppCall.Err
}

func (ms *mockMobileSession) CloseApp() error {
	return ms.CloseAppCall.Err
}

func (ms *mockMobileSession) InstallApp(appPath string) error {
	return ms.InstallAppCall.Err
}

func (ms *mockMobileSession) PerformTouch(actions []interface{}) error {
	return ms.PerformTouchCall.Err
}
