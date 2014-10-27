package session

import "encoding/json"

type Capabilities map[string]interface{}

func (c Capabilities) Browser(browser string) Capabilities {
	c["browserName"] = browser
	return c
}

func (c Capabilities) Version(version string) Capabilities {
	c["version"] = version
	return c
}

func (c Capabilities) Platform(platform string) Capabilities {
	c["platform"] = platform
	return c
}

func (c Capabilities) With(feature string) Capabilities {
	c[feature] = true
	return c
}

func (c Capabilities) Without(feature string) Capabilities {
	c[feature] = false
	return c
}

func (c Capabilities) JSON() string {
	desiredCapabilities := struct {
		DesiredCapabilities map[string]interface{} `json:"desiredCapabilities"`
	}{c}

	json, _ := json.Marshal(desiredCapabilities)

	return string(json)
}
