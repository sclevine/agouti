package types

type Session interface {
	Execute(endpoint, method string, body interface{}, result ...interface{}) error
}
