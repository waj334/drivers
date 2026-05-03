package cyw4343w

const (
	errInvalidCommand         _error = "invalid command"
	errInvalidFirmwareImage   _error = "invalid firmware image"
	errInvalidNvramImage      _error = "invalid nvram image"
	errInvalidClmImage        _error = "invalid clm image"
	errFirmwareDownloadFailed _error = "firmware download failed"
	errCoreIsNotUp            _error = "core is not up"
	errTimeout                _error = "timeout"
	errIoctlFailed            _error = "ioctl failed"
	errJoinFailed             _error = "wifi join failed"
	errPassphraseTooLong      _error = "passphrase exceeds 64 bytes"
	errSSIDTooLong            _error = "SSID exceeds 32 bytes"
)

type _error string

func (e _error) Error() string {
	return string(e)
}
