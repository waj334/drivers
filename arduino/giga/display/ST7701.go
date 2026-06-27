package display

import (
	"time"

	"pkg.si-go.dev/drivers/core/display"
)

const (
	dsiCmd2bkxSel        = 0xFF
	dsiCmd2bk1Sel        = 0x11
	dsiCmd2bk0Sel        = 0x10
	dsiCmd2bkxSelNone    = 0x00
	dsiCmd2Bk0Pvgamctrl  = 0xB0 /* Positive Voltage Gamma Control */
	dsiCmd2Bk0Nvgamctrl  = 0xB1 /* Negative Voltage Gamma Control */
	dsiCmd2Bk0Lneset     = 0xC0 /* Display Line setting */
	dsiCmd2Bk0Porctrl    = 0xC1 /* Porch control */
	dsiCmd2Bk0Invsel     = 0xC2 /* Inversion selection, Frame Rate Control */
	dsiCmd2Bk1Vrhs       = 0xB0 /* Vop amplitude setting */
	dsiCmd2Bk1Vcom       = 0xB1 /* VCOM amplitude setting */
	dsiCmd2Bk1Vghss      = 0xB2 /* VGH Voltage setting */
	dsiCmd2Bk1Testcmd    = 0xB3 /* TEST Command Setting */
	dsiCmd2Bk1Vgls       = 0xB5 /* VGL Voltage setting */
	dsiCmd2Bk1Pwctlr1    = 0xB7 /* Power Control 1 */
	dsiCmd2Bk1Pwctlr2    = 0xB8 /* Power Control 2 */
	dsiCmd2Bk1Spd1       = 0xC1 /* Source pre_drive timing set1 */
	dsiCmd2Bk1Spd2       = 0xC2 /* Source EQ2 Setting */
	dsiCmd2Bk1Mipiset1   = 0xD0 /* MIPI Setting 1 */
	mipiDcsSoftReset     = 0x01
	mipiDcsExitSleepMode = 0x11
)

// initST7701 reproduces Arduino's st7701_init() byte-for-byte, including the
// Display On at the end. All commands are sent via LP packets during video
// blanking on the active DSI stream.
func (g *GigaDisplay) initST7701() error {
	// MIPI_DCS_SOFT_RESET
	if err := g.dcsShortNP(mipiDcsSoftReset); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)

	// MIPI_DCS_EXIT_SLEEP_MODE
	if err := g.dcsShortNP(mipiDcsExitSleepMode); err != nil {
		return err
	}
	time.Sleep(800 * time.Millisecond)

	// ===== Bank0: Display Control =====
	if err := g.genericLongWrite([]byte{0xFF, 0x77, 0x01, 0x00, 0x00, 0x10}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xC0, 0x63, 0x00}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xC1, 0x11, 0x02}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xC2, 0x01, 0x08}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xCC, 0x18}); err != nil {
		return err
	}

	// ===== Gamma =====
	if err := g.genericLongWrite([]byte{
		0xB0,
		0x40, 0xC9, 0x91, 0x0D,
		0x12, 0x07, 0x02, 0x09, 0x09,
		0x1F, 0x04, 0x50, 0x0F, 0xE4,
		0x29, 0xDF,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xB1,
		0x40, 0xCB, 0xD0, 0x11,
		0x92, 0x07, 0x00, 0x08, 0x07,
		0x1C, 0x06, 0x53, 0x12, 0x63,
		0xEB, 0xDF,
	}); err != nil {
		return err
	}

	// ===== Bank1 Select =====
	if err := g.dcsLong([]byte{dsiCmd2bkxSel, 0x77, 0x01, 0x00, 0x00, dsiCmd2bk1Sel}); err != nil {
		return err
	}

	// ===== Bank1: Power Control =====
	if err := g.genShort1P(dsiCmd2Bk1Vrhs, 0x65); err != nil { // VRHS
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Vcom, 0x34); err != nil { // VCOM
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Vghss, 0x87); err != nil { // VGHSS
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Testcmd, 0x80); err != nil { // TESTCMD
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Vgls, 0x49); err != nil { // VGLS
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Pwctlr1, 0x85); err != nil { // PWCTLR1
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Pwctlr2, 0x20); err != nil { // PWCTLR2
		return err
	}
	if err := g.genShort1P(0xB9, 0x10); err != nil {
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Spd1, 0x78); err != nil { // SPD1
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Spd2, 0x78); err != nil { // SPD2
		return err
	}
	if err := g.genShort1P(dsiCmd2Bk1Mipiset1, 0x88); err != nil { // MIPISET1
		return err
	}
	time.Sleep(100 * time.Millisecond)

	// ===== GIP Settings =====
	if err := g.genericLongWrite([]byte{0xE0, 0x00, 0x00, 0x02}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xE1,
		0x08, 0x00, 0x0A, 0x00,
		0x07, 0x00, 0x09, 0x00,
		0x00, 0x33, 0x33,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xE2,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xE3, 0x00, 0x00, 0x33, 0x33}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xE4, 0x44, 0x44}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xE5,
		0x0E, 0x60, 0xA0, 0xA0,
		0x10, 0x60, 0xA0, 0xA0,
		0x0A, 0x60, 0xA0, 0xA0,
		0x0C, 0x60, 0xA0, 0xA0,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xE6, 0x00, 0x00, 0x33, 0x33}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xE7, 0x44, 0x44}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xE8,
		0x0D, 0x60, 0xA0, 0xA0,
		0x0F, 0x60, 0xA0, 0xA0,
		0x09, 0x60, 0xA0, 0xA0,
		0x0B, 0x60, 0xA0, 0xA0,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xEB,
		0x02, 0x01, 0xE4, 0xE4,
		0x44, 0x00, 0x40,
	}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{0xEC, 0x02, 0x01}); err != nil {
		return err
	}
	if err := g.genericLongWrite([]byte{
		0xED,
		0xAB, 0x89, 0x76, 0x54,
		0x01, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0x10,
		0x45, 0x67, 0x98, 0xBA,
	}); err != nil {
		return err
	}

	// ===== Exit command bank =====
	if err := g.genericLongWrite([]byte{dsiCmd2bkxSel, 0x77, 0x01, 0x00, 0x00, dsiCmd2bkxSelNone}); err != nil {
		return err
	}

	time.Sleep(10 * time.Millisecond)

	// COLMOD must match the DSI/panel-side pixel format.
	if err := g.setCOLMODForSurface(); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	// ===== Display On (matches Arduino exactly) =====
	if err := g.dcsShortNP(0x29); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)

	return nil
}

func colmodForPixelFormat(format display.PixelFormat) (byte, error) {
	switch format {
	case display.RGB565:
		return 0x50, nil
	case display.RGB888:
		return 0x70, nil
	default:
		return 0, display.ErrInvalidSurface
	}
}

func (g *GigaDisplay) setCOLMODForSurface() error {
	colmod, err := colmodForPixelFormat(g.surface.Format)
	if err != nil {
		return err
	}

	return g.DSI.DCSShortWriteParam(0, 0x3A, colmod)
}
