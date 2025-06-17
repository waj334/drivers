package cyw4343w

const (
	errInvalidCommand _error = "invalid command"
)

type _error string

func (e _error) Error() string {
	return string(e)
}
