package mocks

type JSON struct {
	ReturnJSON string
}

func (j *JSON) JSON() string {
	return j.ReturnJSON
}
