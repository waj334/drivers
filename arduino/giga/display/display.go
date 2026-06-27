package display

import (
	"fmt"
	"time"

	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/dsihost"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/ltdc"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/pin"
	dsihostreg "pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/dsihost"
	"pkg.si-go.dev/drivers/core/display"
)

const (
	Width  = 480
	Height = 800
)

type GigaDisplay struct {
	LTDC ltdc.LTDC
	DSI  dsihost.DSIHost

	ResetPin     pin.Pin
	TEPin        pin.Pin
	BacklightPin pin.Pin

	fillColor uint32
	surface   display.Surface
	config    PipelineConfig
}

type Config struct {
	PipelineConfig PipelineConfig
	Surface        display.Surface
	FillColor      uint32

	ResetPin     pin.Pin
	TEPin        pin.Pin
	BacklightPin pin.Pin
}

func New(config Config) *GigaDisplay {
	return &GigaDisplay{
		LTDC:         ltdc.LTDC{},
		DSI:          dsihost.DSIHOST,
		ResetPin:     config.ResetPin,
		TEPin:        config.TEPin,
		BacklightPin: config.BacklightPin,

		fillColor: config.FillColor,
		surface:   config.Surface,
		config:    config.PipelineConfig,
	}
}

func (g *GigaDisplay) Bounds() display.Rect {
	return display.Rect{X: 0, Y: 0, W: Width, H: Height}
}

func (g *GigaDisplay) Surface() display.Surface {
	return g.surface
}

func (g *GigaDisplay) SetSurface(surface display.Surface) error {
	if err := validateSurface(surface); err != nil {
		return err
	}

	g.surface = surface
	return g.LTDC.SetLayerFramebuffer(ltdc.Layer1, surface.Ptr, true)
}

func (g *GigaDisplay) Present(surface display.Surface) error {
	if err := validateSurface(surface); err != nil {
		return err
	}

	if err := g.LTDC.SetLayerFramebufferAndWait(
		ltdc.Layer1,
		surface.Ptr,
	); err != nil {
		return err
	}

	// Do not update this until the vertical-blank reload has actually happened.
	g.surface = surface
	return nil
}

// Init mirrors Arduino's dsi_init followed by st7701_init.
//
// The important ordering detail is that LTDC global timing is enabled before
// the DSI host starts, while the LTDC layer registers are configured after
// HAL_DSI_Refresh/WCR.LTDCEN. Arduino does the same: HAL_LTDC_Init,
// HAL_DSI_Start, HAL_DSI_Refresh, then dsi_layerInit/drawCurrentFrameBuffer.
func (g *GigaDisplay) Init() error {
	if err := validateSurface(g.surface); err != nil {
		return err
	}

	// Turn off the backlight during initialization.
	_ = g.SetBacklight(0x00)

	// Fill the buffer with white pixels.
	switch g.surface.Format {
	case display.RGB565:
		fillRGB565(g.surface, uint16(g.fillColor))
	case display.RGB888:
		fillRGB888(g.surface, g.fillColor)
	default:
		return fmt.Errorf("unsupported surface format: %v", g.surface.Format)
	}

	// Hardware reset before the DSI link starts talking to the panel.
	if err := g.Reset(); err != nil {
		_ = g.HardStop()
		return err
	}

	if err := g.configureDSIHost(); err != nil {
		_ = g.HardStop()
		return err
	}

	if err := g.configureLTDC(); err != nil {
		_ = g.HardStop()
		return err
	}

	if err := g.DSI.Start(); err != nil {
		_ = g.HardStop()
		return err
	}

	if err := g.initST7701(); err != nil {
		_ = g.HardStop()
		return err
	}

	g.DSI.EnableLTDC()

	return g.SetBacklight(0xFF)
}

func (g *GigaDisplay) DeInit() error {
	var firstErr error
	remember := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	remember(g.SetBacklight(0))
	remember(g.DisplayOff())
	time.Sleep(20 * time.Millisecond)
	remember(g.SleepIn())

	g.DSI.Stop()
	g.LTDC.Disable()

	if g.ResetPin != pin.NoPin {
		g.ResetPin.Low()
	}

	return firstErr
}

func (g *GigaDisplay) HardStop() error {
	_ = g.SetBacklight(0)

	g.DSI.Stop()
	g.LTDC.Disable()

	if g.ResetPin != pin.NoPin {
		g.ResetPin.Low()
	}
	return nil
}

func (g *GigaDisplay) Reset() error {
	if g.ResetPin == pin.NoPin {
		return nil
	}

	g.ResetPin.SetOutputMode(pin.PushPull)
	g.ResetPin.SetPullMode(pin.NoPull)
	g.ResetPin.SetSpeedMode(pin.LowSpeed)
	g.ResetPin.SetMode(pin.Output)

	g.ResetPin.Low()
	time.Sleep(10 * time.Millisecond)

	g.ResetPin.High()
	time.Sleep(120 * time.Millisecond)
	return nil
}

func (g *GigaDisplay) SleepOut() error {
	if err := g.dcsShortNP(0x11); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)
	return nil
}

func (g *GigaDisplay) SleepIn() error {
	return g.dcsShortNP(0x10)
}

func (g *GigaDisplay) DisplayOn() error {
	return g.dcsShortNP(0x29)
}

func (g *GigaDisplay) DisplayOff() error {
	return g.dcsShortNP(0x28)
}

func (g *GigaDisplay) SetBacklight(percent uint8) error {
	if g.BacklightPin == pin.NoPin {
		return nil
	}

	g.BacklightPin.SetOutputMode(pin.PushPull)
	g.BacklightPin.SetPullMode(pin.NoPull)
	g.BacklightPin.SetSpeedMode(pin.LowSpeed)
	g.BacklightPin.SetMode(pin.Output)

	g.BacklightPin.Set(percent != 0)
	return nil
}

func (g *GigaDisplay) configureLTDC() error {
	return g.LTDC.Configure(g.config.LTDC)
}

func (g *GigaDisplay) configureDSIHost() error {
	cfg := g.config.DSI
	cfg.Enable = true
	cfg.LPCommandMode = true
	return g.DSI.Configure(cfg)
}

func validateSurface(s display.Surface) error {
	if s.Ptr == nil {
		return display.ErrInvalidSurface
	}
	if s.Width != Width || s.Height != Height {
		return display.ErrInvalidSurface
	}

	bpp, err := bytesPerPixel(s.Format)
	if err != nil {
		return err
	}

	minStride := Width * bpp
	minLen := minStride * Height

	if s.Stride < minStride {
		return display.ErrInvalidSurface
	}
	if s.Len < uintptr(minLen) {
		return display.ErrInvalidSurface
	}
	return nil
}

// dcsLong wraps DSI.DCSLongWrite with diagnostic dumps and a brief inter-command delay.
func (g *GigaDisplay) dcsLong(packet []byte) error {
	if len(packet) == 0 {
		return nil
	}

	for dsihostreg.Dsihost.Dsigpsr.GetCmdff() {
	}

	err := g.DSI.DCSLongWrite(0, packet[0], packet[1:])
	if err != nil {
		return err
	}

	return nil
}

func (g *GigaDisplay) genericLongWrite(packet []byte) error {
	if len(packet) == 0 {
		return nil
	}

	const maxParams = 19
	if len(packet)-1 > maxParams {
		return display.ErrInvalidConfig
	}

	for dsihostreg.Dsihost.Dsigpsr.GetCmdff() {
	}

	err := g.DSI.DCSLongWrite(0, packet[0], packet[1:])
	if err != nil {
		return err
	}

	return nil
}

// genShort1P sends DT=0x13 with two data bytes.
func (g *GigaDisplay) genShort1P(cmd byte, param byte) error {
	for dsihostreg.Dsihost.Dsigpsr.GetCmdff() {
	}

	err := g.DSI.GenericShortWrite1P(0, cmd, param)
	if err != nil {
		return err
	}
	return nil
}

func (g *GigaDisplay) dcsShortNP(cmd byte) error {
	for dsihostreg.Dsihost.Dsigpsr.GetCmdff() {
	}

	if err := g.DSI.DCSShortWrite(0, cmd); err != nil {
		return err
	}
	return nil
}

func (g *GigaDisplay) dcsRead(cmd byte) (byte, error) {
	var buf [1]byte
	n, err := g.DSI.DCSRead(0, cmd, buf[:])
	if err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, display.ErrInvalidConfig
	}
	return buf[0], nil
}

func bytesPerPixel(format display.PixelFormat) (int, error) {
	switch format {
	case display.RGB565:
		return 2, nil
	case display.RGB888:
		return 3, nil
	default:
		return 0, display.ErrInvalidSurface
	}
}

func ltdcPixelFormat(format display.PixelFormat) (ltdc.PixelFormat, error) {
	switch format {
	case display.RGB565:
		return ltdc.PixelFormatRGB565, nil
	case display.RGB888:
		return ltdc.PixelFormatRGB888, nil
	default:
		return 0, display.ErrInvalidSurface
	}
}
