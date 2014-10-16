package types

type Cookie struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	HTTPOnly bool        `json:"httpOnly"`
	Expiry   int64       `json:"expiry"`
}
