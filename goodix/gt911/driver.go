package gt911

import (
	"time"

	"pkg.si-go.dev/chip/core/hal/i2c"
	"pkg.si-go.dev/chip/core/hal/pin"
)

const (
	// GT911 can be selected as either 0x5D/0xBA or 0x14/0x28.
	// Keep the 8-bit write-address default because your existing skeleton used 0xBA.
	// If your I2C HAL expects 7-bit addresses, set Config.Address to Address5D or Address14.
	AddressBA uint16 = 0xBA
	Address28 uint16 = 0x28
	Address5D uint16 = 0x5D
	Address14 uint16 = 0x14

	busAddress = AddressBA

	MaxI2CClockHz  = 400_000
	MaxTouchPoints = 5
)

const (
	regCommand        = 0x8040
	regESDCheck       = 0x8041
	regCommandCheck   = 0x8046
	regConfigStart    = 0x8047
	regConfigVersion  = 0x8047
	regXOutputMax     = 0x8048
	regYOutputMax     = 0x804A
	regTouchNumber    = 0x804C
	regModuleSwitch1  = 0x804D
	regConfigEnd      = 0x80FE
	regConfigChecksum = 0x80FF
	regConfigFresh    = 0x8100

	regProductID        = 0x8140
	regCoordinateStatus = 0x814E
	regTouchDataStart   = 0x814F
)

const (
	commandReadCoordinates = 0x00
	commandReadDiffData    = 0x01
	commandReadRawData     = 0x02
	commandBaselineUpdate  = 0x03
	commandCalibrate       = 0x04
	commandScreenOff       = 0x05
	commandEnterChargeMode = 0x06
	commandExitChargeMode  = 0x07
	commandGestureMode     = 0x08

	commandESDCheck = 0xAA
)

const (
	coordinateStatusReady       = 0x80
	coordinateStatusLargeDetect = 0x40
	coordinateStatusHaveKey     = 0x10
	coordinateStatusTouchMask   = 0x0F

	touchRecordSize  = 8
	configDataLength = regConfigEnd - regConfigStart + 1
	maxWriteLength   = configTotalLength + 2
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrShortWrite      Error = "gt911: short i2c write"
	ErrShortRead       Error = "gt911: short i2c read"
	ErrWriteTooLarge   Error = "gt911: write is too large for driver scratch buffer"
	ErrBadTouchCount   Error = "gt911: invalid touch count"
	ErrBufferTooSmall  Error = "gt911: touch buffer is too small"
	ErrBadConfigLength Error = "gt911: bad configuration length"
	ErrClockFrequency  Error = "gt911: unsupported i2c clock frequency"
)

type TouchPoint struct {
	Valid bool
	ID    uint8
	X     uint16
	Y     uint16
	Size  uint16
}

type TouchStatus struct {
	Ready       bool
	LargeDetect bool
	HaveKey     bool
	Count       uint8
}

type Info struct {
	ProductID   [4]byte
	Firmware    uint16
	XResolution uint16
	YResolution uint16
	VendorID    uint8
}

type InterruptTrigger uint8

const (
	InterruptRisingEdge  InterruptTrigger = 0
	InterruptFallingEdge InterruptTrigger = 1
	InterruptLowLevel    InterruptTrigger = 2
	InterruptHighLevel   InterruptTrigger = 3
)

type TouchHandler func([]TouchPoint)

type GT911[
	TransportT i2c.I2C,
	PinT pin.Pin[PinT],
] struct {
	i2c   TransportT
	irq   PinT
	reset PinT

	address          uint16
	clockFrequencyHz uint32

	xResolution uint16
	yResolution uint16
	maxTouches  uint8

	inputMode  pin.Mode
	outputMode pin.Mode
	irqMode    pin.IRQMode

	onTouch         TouchHandler
	readInInterrupt bool
	irqPending      bool
	lastErr         error

	info       Info
	touchCount int
	touches    [MaxTouchPoints]TouchPoint

	// One shared scratch buffer keeps register writes allocation-free.
	scratch [maxWriteLength]byte
}

type Config[
	TransportT i2c.I2C,
	PinT pin.Pin[PinT],
] struct {
	I2C   TransportT
	IRQ   PinT
	Reset PinT

	// Address is the value expected by the target I2C HAL. Your current skeleton used
	// 0xBA. Use 0x5D instead if the HAL expects 7-bit addresses.
	Address uint16

	// Optional. If non-zero, Init sets the I2C bus to this frequency.
	ClockFrequencyHz uint32

	XResolution uint16
	YResolution uint16
	MaxTouches  uint8

	InputMode  pin.Mode
	OutputMode pin.Mode
	IRQMode    pin.IRQMode

	// OnTouch is called by HandleInterrupt/Poll after a successful coordinate read.
	// The slice is backed by the driver; copy it if it must outlive the callback.
	OnTouch TouchHandler

	// Leave this false for normal ISR use. When false, the IRQ handler only marks
	// the driver pending and clears the MCU interrupt; call HandleInterrupt outside
	// the ISR to perform the I2C transaction.
	ReadInInterrupt bool
}

func NewGT911[
	TransportT i2c.I2C,
	PinT pin.Pin[PinT],
](config Config[TransportT, PinT]) (*GT911[TransportT, PinT], error) {
	address := config.Address
	if address == 0 {
		address = busAddress
	}

	maxTouches := config.MaxTouches
	if maxTouches == 0 {
		maxTouches = MaxTouchPoints
	}
	if maxTouches > MaxTouchPoints {
		maxTouches = MaxTouchPoints
	}

	return &GT911[TransportT, PinT]{
		i2c:   config.I2C,
		irq:   config.IRQ,
		reset: config.Reset,

		address:          address,
		clockFrequencyHz: config.ClockFrequencyHz,

		xResolution: config.XResolution,
		yResolution: config.YResolution,
		maxTouches:  maxTouches,

		inputMode:  config.InputMode,
		outputMode: config.OutputMode,
		irqMode:    config.IRQMode,

		onTouch:         config.OnTouch,
		readInInterrupt: config.ReadInInterrupt,
	}, nil
}

func (g *GT911[TransportT, PinT]) Init() error {
	if g.clockFrequencyHz != 0 {
		clock := g.clockFrequencyHz
		if clock > MaxI2CClockHz {
			clock = MaxI2CClockHz
		}
		if !g.i2c.SetClockFrequency(clock) {
			return ErrClockFrequency
		}
	}

	g.Reset()

	if err := g.writeCommand(commandReadCoordinates); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	info, err := g.ReadInfo()
	if err != nil {
		return err
	}
	g.info = info

	if err := g.SetupConfig(g.maxTouches); err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	if err := g.ClearTouchStatus(); err != nil {
		return err
	}

	g.irq.SetInterrupt(g.irqMode, g.irqHandler)
	return nil
}

func (g *GT911[TransportT, PinT]) SetupConfig(maxTouches uint8) error {
	var cfg [configTotalLength]byte

	if err := g.readRegisters(regConfigStart, cfg[:]); err != nil {
		return err
	}

	if cfg[0] == 0 {
		//return ErrBadConfigVersion
	}

	if ConfigChecksum(cfg[:configPayloadLength]) != cfg[configPayloadLength] {
		return ErrBadConfigChecksum
	}

	if maxTouches == 0 {
		maxTouches = MaxTouchPoints
	}
	if maxTouches > MaxTouchPoints {
		maxTouches = MaxTouchPoints
	}

	cfg[configTouchNumberOffset] = maxTouches
	cfg[configPayloadLength] = ConfigChecksum(cfg[:configPayloadLength])
	cfg[configPayloadLength+1] = 1 // Config_Fresh

	return g.writeRegisters(regConfigStart, cfg[:])
}

func (g *GT911[TransportT, PinT]) Reset() {
	g.reset.SetMode(g.outputMode)

	g.irq.SetInterrupt(g.irqMode, nil)

	// INT selects address during reset.
	g.irq.SetMode(g.outputMode)
	g.irq.Set(addressSelectLevel(g.address))

	// Conservative power/reset settle.
	time.Sleep(20 * time.Millisecond)

	// Active-low reset.
	g.reset.Set(false)
	time.Sleep(10 * time.Millisecond)

	g.reset.Set(true)

	// Keep INT driven during the post-reset address-select window.
	time.Sleep(50 * time.Millisecond)

	// Now release INT as input/floating.
	g.irq.SetMode(g.inputMode)

	// Let firmware finish booting before I2C/config.
	time.Sleep(50 * time.Millisecond)
}

func (g *GT911[TransportT, PinT]) Info() Info {
	return g.info
}

func (g *GT911[TransportT, PinT]) LastError() error {
	return g.lastErr
}

func (g *GT911[TransportT, PinT]) Pending() bool {
	return g.irqPending
}

func (g *GT911[TransportT, PinT]) Touches() []TouchPoint {
	return g.touches[:g.touchCount]
}

func (g *GT911[TransportT, PinT]) ReadInfo() (Info, error) {
	var buf [11]byte
	if err := g.readRegisters(regProductID, buf[:]); err != nil {
		return Info{}, err
	}

	return Info{
		ProductID:   [4]byte{buf[0], buf[1], buf[2], buf[3]},
		Firmware:    u16(buf[4], buf[5]),
		XResolution: u16(buf[6], buf[7]),
		YResolution: u16(buf[8], buf[9]),
		VendorID:    buf[10],
	}, nil
}

func (g *GT911[TransportT, PinT]) ReadTouchStatus() (TouchStatus, error) {
	var buf [1]byte
	if err := g.readRegisters(regCoordinateStatus, buf[:]); err != nil {
		return TouchStatus{}, err
	}
	return decodeStatus(buf[0]), nil
}

func (g *GT911[TransportT, PinT]) ReadTouchPoints(dst []TouchPoint) (TouchStatus, error) {
	status, err := g.ReadTouchStatus()
	if err != nil || !status.Ready {
		return status, err
	}

	count := int(status.Count)
	if count > MaxTouchPoints {
		_ = g.ClearTouchStatus()
		return status, ErrBadTouchCount
	}
	if count > len(dst) {
		return status, ErrBufferTooSmall
	}
	if count == 0 {
		return status, g.ClearTouchStatus()
	}

	dataLen := count * touchRecordSize
	buf := g.scratch[:dataLen]
	if err := g.readRegisters(regTouchDataStart, buf); err != nil {
		return status, err
	}

	for i := range len(dst) {
		if i < count {
			base := i * touchRecordSize
			dst[i] = TouchPoint{
				Valid: true,
				ID:    buf[base+0],
				X:     u16(buf[base+1], buf[base+2]),
				Y:     u16(buf[base+3], buf[base+4]),
				Size:  u16(buf[base+5], buf[base+6]),
			}
		} else {
			dst[i].Valid = false
		}
	}

	if err := g.ClearTouchStatus(); err != nil {
		return status, err
	}
	return status, nil
}

func (g *GT911[TransportT, PinT]) Poll() (TouchStatus, []TouchPoint, error) {
	status, err := g.ReadTouchPoints(g.touches[:])
	if err != nil {
		g.lastErr = err
		return status, nil, err
	}

	if status.Ready {
		g.touchCount = int(status.Count)
	} else {
		g.touchCount = 0
	}

	touches := g.touches[:g.touchCount]
	if status.Ready && g.onTouch != nil {
		g.onTouch(touches)
	}
	return status, touches, nil
}

func (g *GT911[TransportT, PinT]) HandleInterrupt() (TouchStatus, []TouchPoint, error) {
	g.irqPending = false
	return g.Poll()
}

func (g *GT911[TransportT, PinT]) ClearTouchStatus() error {
	return g.writeRegister(regCoordinateStatus, 0)
}

func (g *GT911[TransportT, PinT]) Sleep() error {
	// GT911 expects INT low before the screen-off command.
	g.irq.SetInterrupt(g.irqMode, nil)
	g.irq.SetMode(g.outputMode)
	g.irq.Set(false)
	time.Sleep(5 * time.Millisecond)

	if err := g.writeCommand(commandScreenOff); err != nil {
		return err
	}

	// The guide requires >58ms between screen-off and wake.
	time.Sleep(60 * time.Millisecond)
	return nil
}

func (g *GT911[TransportT, PinT]) Wake() {
	// Wake from sleep by driving INT high for 2ms to 5ms.
	g.irq.SetMode(g.outputMode)
	g.irq.Set(true)
	time.Sleep(3 * time.Millisecond)
	g.irq.SetMode(g.inputMode)
	time.Sleep(50 * time.Millisecond)
	g.irq.SetInterrupt(g.irqMode, g.irqHandler)
}

func (g *GT911[TransportT, PinT]) EnterGestureMode() error {
	return g.writeCommand(commandGestureMode)
}

func (g *GT911[TransportT, PinT]) EnterChargeMode() error {
	return g.writeCommand(commandEnterChargeMode)
}

func (g *GT911[TransportT, PinT]) ExitChargeMode() error {
	return g.writeCommand(commandExitChargeMode)
}

func (g *GT911[TransportT, PinT]) EnableESDCheck() error {
	return g.writeRegister(regESDCheck, commandESDCheck)
}

func (g *GT911[TransportT, PinT]) ReadConfig(dst []byte) error {
	if len(dst) != configDataLength {
		return ErrBadConfigLength
	}
	return g.readRegisters(regConfigStart, dst)
}

func (g *GT911[TransportT, PinT]) WriteConfig(config []byte) error {
	if len(config) != configDataLength {
		return ErrBadConfigLength
	}
	if err := g.writeRegisters(regConfigStart, config); err != nil {
		return err
	}

	fresh := [2]byte{ConfigChecksum(config), 1}
	return g.writeRegisters(regConfigChecksum, fresh[:])
}

func (g *GT911[TransportT, PinT]) SetResolution(x, y uint16, touches uint8) error {
	var config [configDataLength]byte
	if err := g.ReadConfig(config[:]); err != nil {
		return err
	}

	config[regXOutputMax-regConfigStart+0] = byte(x)
	config[regXOutputMax-regConfigStart+1] = byte(x >> 8)
	config[regYOutputMax-regConfigStart+0] = byte(y)
	config[regYOutputMax-regConfigStart+1] = byte(y >> 8)
	if touches != 0 {
		if touches > MaxTouchPoints {
			touches = MaxTouchPoints
		}
		config[regTouchNumber-regConfigStart] = touches
	}

	return g.WriteConfig(config[:])
}

func (g *GT911[TransportT, PinT]) SetInterruptTrigger(trigger InterruptTrigger) error {
	var config [configDataLength]byte
	if err := g.ReadConfig(config[:]); err != nil {
		return err
	}

	idx := regModuleSwitch1 - regConfigStart
	config[idx] = (config[idx] &^ 0x03) | (byte(trigger) & 0x03)
	return g.WriteConfig(config[:])
}

func ConfigChecksum(config []byte) byte {
	var sum uint8
	for _, b := range config {
		sum += b
	}
	return ^sum + 1
}

func (g *GT911[TransportT, PinT]) irqHandler(p PinT) {
	g.irqPending = true

	if !g.readInInterrupt {
		return
	}

	_, _, err := g.HandleInterrupt()
	if err != nil {
		g.lastErr = err
	}
}

func (g *GT911[TransportT, PinT]) writeCommand(command byte) error {
	if command > commandExitChargeMode {
		if err := g.writeRegister(regCommandCheck, command); err != nil {
			return err
		}
	}
	return g.writeRegister(regCommand, command)
}

func (g *GT911[TransportT, PinT]) writeRegister(reg uint16, value byte) error {
	buf := [1]byte{value}
	return g.writeRegisters(reg, buf[:])
}

func (g *GT911[TransportT, PinT]) writeRegisters(reg uint16, data []byte) error {
	if len(data)+2 > len(g.scratch) {
		return ErrWriteTooLarge
	}

	buf := g.scratch[:len(data)+2]
	buf[0] = byte(reg >> 8)
	buf[1] = byte(reg)
	copy(buf[2:], data)

	n, err := g.i2c.WriteAddress(g.address, buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return ErrShortWrite
	}
	return nil
}

func (g *GT911[TransportT, PinT]) readRegisters(reg uint16, dst []byte) error {
	addr := [2]byte{byte(reg >> 8), byte(reg)}

	n, err := g.i2c.WriteAddress(g.address, addr[:])
	if err != nil {
		return err
	}
	if n != len(addr) {
		return ErrShortWrite
	}

	n, err = g.i2c.ReadAddress(g.address, dst)
	if err != nil {
		return err
	}
	if n != len(dst) {
		return ErrShortRead
	}
	return nil
}

func decodeStatus(value byte) TouchStatus {
	return TouchStatus{
		Ready:       value&coordinateStatusReady != 0,
		LargeDetect: value&coordinateStatusLargeDetect != 0,
		HaveKey:     value&coordinateStatusHaveKey != 0,
		Count:       value & coordinateStatusTouchMask,
	}
}

func addressSelectLevel(address uint16) bool {
	switch address {
	case Address28, Address14:
		return true
	default:
		return false
	}
}

func u16(lo, hi byte) uint16 {
	return uint16(lo) | uint16(hi)<<8
}
