package mobile

import (
	"github.com/sclevine/agouti/api"
)

type Session struct {
	*api.Session
}

// override methods like this and return mobile.Element!
func (s *Session) GetElement(selector api.Selector) (*Element, error) {
	apiElement, err := s.Session.GetElement(selector)
	if err != nil {
		return nil, err
	}

	return &Element{apiElement, s}, nil
}

//
// Appium-centric functions
//

func (s *Session) PerformTouch(actions []interface{}) error {
	request := struct {
		Actions []interface{} `json:"actions"`
	}{actions}

	return s.Send("POST", "touch/perform", request, nil)
}


func (s *Session) InstallApp(appPath string) error {
	request := struct {
		AppPath string `json:"appPath"`
	}{appPath}

	return s.Send("appium/device/install_app", "POST", request, nil)
}

func (s *Session) RemoveApp(appId string) error {
	request := struct {
		AppID string `json:"appId"`
	}{appId}

	return s.Send("appium/device/remove_app", "POST", request, nil)
}

func (s *Session) IsAppInstalled(bundleId string) (bool, error) {
	request := struct {
		BundleID string `json:"bundleId"`
	}{bundleId}

	var out bool
	if err := s.Send("appium/device/app_installed", "POST", request, &out); err != nil {
		return false, err
	}
	return out, nil
}

func (s *Session) LaunchApp() error {
	return s.Send("appium/app/launch", "POST", nil, nil)
}

func (s *Session) CloseApp() error {
	return s.Send("appium/app/launch", "POST", nil, nil)
}

func (s *Session) GetAppStrings(language string) ([]string, error) {
	request := struct {
		Language string `json:"language"`
	}{language}

	var strs []string
	if err := s.Send("appium/app/strings", "POST", request, &strs); err != nil {
		return nil, err
	}
	return strs, nil
}

func (s *Session) GetCurrentActivity() (string, error) {
	var activity string
	if err := s.Send("appium/device/current_activity", "GET", nil, &activity); err != nil {
		return "", err
	}
	return activity, nil
}

func (s *Session) Lock() error {
	return s.Send("appium/device/lock", "POST", nil, nil)
}

func (s *Session) Shake() error {
	return s.Send("appium/device/shake", "POST", nil, nil)
}

func (s *Session) Reset() error {
	return s.Send("appium/app/reset", "POST", nil, nil)
}

func (s *Session) OpenNotifications() error {
	return s.Send("appium/device/open_notifications", "POST", nil, nil)
}

func (s *Session) GetSettings() (map[string]interface{}, error) {
	var out map[string]interface{}

	if err := s.Send("appium/settings", "GET", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Session) UpdateSettings(settings map[string]interface{}) error {
	request := struct {
		Settings map[string]interface{} `json:"settings"`
	}{settings}

	return s.Send("appium/settings", "POST", request, nil)
}

func (s *Session) ToggleLocationServices() error {
	return s.Send("appium/device/toggle_location_services", "POST", nil, nil)
}

var _ = `
        self.command_executor._commands[Command.CONTEXTS] = \
            ('GET', '/session/$sessionId/contexts')
        self.command_executor._commands[Command.GET_CURRENT_CONTEXT] = \
            ('GET', '/session/$sessionId/context')
        self.command_executor._commands[Command.SWITCH_TO_CONTEXT] = \
            ('POST', '/session/$sessionId/context')
        self.command_executor._commands[Command.TOUCH_ACTION] = \
            ('POST', '/session/$sessionId/touch/perform')
        self.command_executor._commands[Command.MULTI_ACTION] = \
            ('POST', '/session/$sessionId/touch/multi/perform')
        self.command_executor._commands[Command.GET_APP_STRINGS] = \
            ('POST', '/session/$sessionId/appium/app/strings')
        # Needed for Selendroid
        self.command_executor._commands[Command.KEY_EVENT] = \
            ('POST', '/session/$sessionId/appium/device/keyevent')
        self.command_executor._commands[Command.PRESS_KEYCODE] = \
            ('POST', '/session/$sessionId/appium/device/press_keycode')
        self.command_executor._commands[Command.LONG_PRESS_KEYCODE] = \
            ('POST', '/session/$sessionId/appium/device/long_press_keycode')
        self.command_executor._commands[Command.GET_CURRENT_ACTIVITY] = \
            ('GET', '/session/$sessionId/appium/device/current_activity')
        self.command_executor._commands[Command.SET_IMMEDIATE_VALUE] = \
            ('POST', '/session/$sessionId/appium/element/$elementId/value')
        self.command_executor._commands[Command.PULL_FILE] = \
            ('POST', '/session/$sessionId/appium/device/pull_file')
        self.command_executor._commands[Command.PULL_FOLDER] = \
            ('POST', '/session/$sessionId/appium/device/pull_folder')
        self.command_executor._commands[Command.PUSH_FILE] = \
            ('POST', '/session/$sessionId/appium/device/push_file')
        self.command_executor._commands[Command.BACKGROUND] = \
            ('POST', '/session/$sessionId/appium/app/background')
        self.command_executor._commands[Command.IS_APP_INSTALLED] = \
            ('POST', '/session/$sessionId/appium/device/app_installed')
        self.command_executor._commands[Command.INSTALL_APP] = \
            ('POST', '/session/$sessionId/appium/device/install_app')
        self.command_executor._commands[Command.REMOVE_APP] = \
            ('POST', '/session/$sessionId/appium/device/remove_app')
        self.command_executor._commands[Command.START_ACTIVITY] = \
            ('POST', '/session/$sessionId/appium/device/start_activity')
        self.command_executor._commands[Command.LAUNCH_APP] = \
            ('POST', '/session/$sessionId/appium/app/launch')
        self.command_executor._commands[Command.CLOSE_APP] = \
            ('POST', '/session/$sessionId/appium/app/close')
        self.command_executor._commands[Command.END_TEST_COVERAGE] = \
            ('POST', '/session/$sessionId/appium/app/end_test_coverage')
        self.command_executor._commands[Command.LOCK] = \
           ('POST', '/session/$sessionId/appium/device/lock')
        self.command_executor._commands[Command.SHAKE] = \
            ('POST', '/session/$sessionId/appium/device/shake')
        self.command_executor._commands[Command.RESET] = \
            ('POST', '/session/$sessionId/appium/app/reset')
        self.command_executor._commands[Command.HIDE_KEYBOARD] = \
            ('POST', '/session/$sessionId/appium/device/hide_keyboard')
        self.command_executor._commands[Command.OPEN_NOTIFICATIONS] = \
            ('POST', '/session/$sessionId/appium/device/open_notifications')
        self.command_executor._commands[Command.GET_NETWORK_CONNECTION] = \
            ('GET', '/session/$sessionId/network_connection')
        self.command_executor._commands[Command.SET_NETWORK_CONNECTION] = \
            ('POST', '/session/$sessionId/network_connection')
        self.command_executor._commands[Command.GET_AVAILABLE_IME_ENGINES] = \
            ('GET', '/session/$sessionId/ime/available_engines')
        self.command_executor._commands[Command.IS_IME_ACTIVE] = \
            ('GET', '/session/$sessionId/ime/activated')
        self.command_executor._commands[Command.ACTIVATE_IME_ENGINE] = \
            ('POST', '/session/$sessionId/ime/activate')
        self.command_executor._commands[Command.DEACTIVATE_IME_ENGINE] = \
            ('POST', '/session/$sessionId/ime/deactivate')
        self.command_executor._commands[Command.GET_ACTIVE_IME_ENGINE] = \
            ('GET', '/session/$sessionId/ime/active_engine')
        self.command_executor._commands[Command.REPLACE_KEYS] = \
            ('POST', '/session/$sessionId/appium/element/$elementId/replace_value')
        self.command_executor._commands[Command.GET_SETTINGS] = \
            ('GET', '/session/$sessionId/appium/settings')
        self.command_executor._commands[Command.UPDATE_SETTINGS] = \
            ('POST', '/session/$sessionId/appium/settings')
        self.command_executor._commands[Command.TOGGLE_LOCATION_SERVICES] = \
            ('POST', '/session/$sessionId/appium/device/toggle_location_services')

:
`
