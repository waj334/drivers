package gt911

const (
	configPayloadLength = regConfigEnd - regConfigStart + 1 // 184: 0x8047..0x80FE
	configTotalLength   = configPayloadLength + 2           // + checksum + fresh = 186

	configTouchNumberOffset = regTouchNumber - regConfigStart
)

const (
	ErrBadConfigChecksum Error = "gt911: bad config checksum"
	ErrBadConfigVersion  Error = "gt911: bad config version"
)

func (g *GT911[TransportT, PinT]) VerifyConfig() error {
	var cfg [configTotalLength]byte

	if err := g.readRegisters(regConfigStart, cfg[:]); err != nil {
		return err
	}

	if cfg[0] == 0 {
		//return ErrBadConfigVersion
	}

	want := ConfigChecksum(cfg[:configPayloadLength])
	got := cfg[configPayloadLength]

	if want != got {
		return ErrBadConfigChecksum
	}

	return nil
}

func (g *GT911[TransportT, PinT]) ReadInterruptTrigger() (InterruptTrigger, error) {
	var b [1]byte
	if err := g.readRegisters(regModuleSwitch1, b[:]); err != nil {
		return 0, err
	}
	return InterruptTrigger(b[0] & 0x03), nil
}
