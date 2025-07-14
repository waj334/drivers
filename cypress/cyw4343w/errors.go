package cyw4343w

const (
	errInvalidCommand         _error = "invalid command"
	errInvalidFirmwareImage   _error = "invalid firmware image"
	errInvalidNvramImage      _error = "invalid nvram image"
	errFirmwareDownloadFailed _error = "firmware download failed"
	errCoreIsNotUp            _error = "core is not up"
)

type _error string

func (e _error) Error() string {
	return string(e)
}
