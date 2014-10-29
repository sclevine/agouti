package types

type Service interface {
	Start() error
	Stop()
	CreateSession(capabilities JSON) (Session, error)
}
