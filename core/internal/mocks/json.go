package mocks

type JSON struct {
	ReturnJSON string
	Err        error
}

func (j *JSON) JSON() (string, error) {
	return j.ReturnJSON, j.Err
}
