package display

import (
	"fmt"

	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/dsihost"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/ltdc"
	"pkg.si-go.dev/drivers/core/display"
)

// Panel timing for the ARDUINO_GIGA + 480x800 IVO panel, matching Arduino's
// envie_known_modes[EDID_MODE_480x800_60Hz]. These are fixed properties of the
// panel, not of the framebuffer surface.
//
// Clock plan (with HSE=16 MHz):
//
//	PLL: IDF=/4, NDIV=125, ODF=/1 -> PHY = 16/4*125/1 = 500 MHz
//	Lane byte clock = 500/8 = 62.5 MHz
//	TX escape clock = 62.5/4 = 15.625 MHz
//
// The LTDC pixel clock (PLL3 in clocks.go) targets 38 MHz. The horizontal DSI
// timings are scaled from the pixel-clock domain into lane-byte cycles
// (Arduino's LANE_BYTE_CLOCK/pixel_clock ratio):
//
//	HSA	= 24 * 62500 / 38000 = 39
//	HBP	= 30 * 62500 / 38000 = 49
//	line = 854 * 62500 / 38000 = 1404
const (
	hSyncLen    = 24
	hBackPorch  = 30
	hFrontPorch = 320

	vSyncLen    = 4
	vBackPorch  = 50
	vFrontPorch = 20
)

// dsiColorFormat maps a surface pixel format to the matching DSI host color
// coding (LCOLCR.COLC) and wrapper color multiplexing (WCFGR.COLMUX). The two
// MUST stay in lockstep: COLMUX is what tells the LTDC->DSI bridge how many
// bytes per pixel to pull out of the LTDC FIFO, so a mismatch against COLC
// makes the link under/over-fetch and the active line is truncated or wrapped.
func dsiColorFormat(format display.PixelFormat) (dsihost.ColorCoding, dsihost.ColorMux, error) {
	switch format {
	case display.RGB565:
		return dsihost.ColorRGB565, dsihost.ColorMux16Bit1, nil

	case display.RGB888, display.XRGB8888, display.ARGB8888:
		return dsihost.ColorRGB888, dsihost.ColorMux24Bit, nil

	default:
		return 0, 0, fmt.Errorf("display: format %v is not supported on the DSI link", format)
	}
}

// DSIConfigForSurface builds the DSI host configuration for the given surface.
// Color coding/packing comes from surface.Format; the active region comes from
// surface dimensions. Porch/sync timings are panel-fixed, so the surface must
// match the panel's native resolution.
func DSIConfigForSurface(surface display.Surface, laneByteClockHz, pixelClockHz uint32) (dsihost.Config, error) {
	if surface.Width != Width || surface.Height != Height {
		return dsihost.Config{}, fmt.Errorf(
			"display: surface %dx%d does not match panel %dx%d",
			surface.Width, surface.Height, Width, Height,
		)
	}

	colorCoding, colorMux, err := dsiColorFormat(surface.Format)
	if err != nil {
		return dsihost.Config{}, err
	}

	activeW := uint32(surface.Width)
	activeH := uint32(surface.Height)

	scale := func(pixels uint32) uint32 {
		return uint32(uint64(pixels) * uint64(laneByteClockHz) / uint64(pixelClockHz))
	}

	hsa := scale(hSyncLen)
	hbp := scale(hBackPorch)
	// HLINE is the TOTAL line in the pixel domain, scaled once to lane-byte cycles.
	hline := scale(hSyncLen + hBackPorch + activeW + hFrontPorch)

	return dsihost.Config{
		Enable: true,

		PLL: dsihost.PLLConfig{
			Enable:          true,
			NDIV:            125,
			IDF:             dsihost.PLLIDF4,
			ODF:             dsihost.PLLODF1,
			RegulatorEnable: true,
		},

		Lanes:         dsihost.TwoDataLanes,
		TXEscapeCkdiv: 4,

		AutomaticClockLaneControl: false,

		VirtualChannelID: 0,
		ColorCoding:      colorCoding,
		ColorMux:         colorMux,
		LooselyPacked:    false,

		// Panel hpol=1, vpol=1. DE active-high on the DSI side; the LTDC side is
		// inverted by the bridge (see display.go configureLTDC notes).
		HSPolarityHigh: true,
		VSPolarityHigh: true,
		DEPolarityHigh: true,

		VideoMode: dsihost.VideoModeBurst,

		// PacketSize (VPSIZE) is the active PIXEL count, never a byte count. The
		// DSI derives bytes/pixel from ColorCoding/ColorMux; sizing this in bytes
		// is what produced the half-width (480->240px) and wrapped-line artifacts.
		PacketSize:     uint16(activeW),
		NumberOfChunks: 0,
		NullPacketSize: 0x0FFF,

		HorizontalSyncActive: uint16(hsa),
		HorizontalBackPorch:  uint16(hbp),
		HorizontalLine:       uint16(hline),

		VerticalSyncActive: vSyncLen,
		VerticalBackPorch:  vBackPorch,
		VerticalFrontPorch: vFrontPorch,
		VerticalActive:     uint16(activeH),

		LPCommandEnable:         true,
		LPLargestPacketSize:     24,
		LPVACTLargestPacketSize: 0,

		LPHFPEnable:  true,
		LPHBPEnable:  true,
		LPVACTEnable: true,
		LPVFPEnable:  true,
		LPVBPEnable:  true,
		LPVSAEnable:  true,

		PhyTiming: dsihost.PhyTiming{
			Enable: true,

			ClockHS2LP: 35,
			ClockLP2HS: 35,

			DataHS2LP: 35,
			DataLP2HS: 35,

			MaxReadTime:  10,
			StopWaitTime: 10,
		},

		LPCommandMode: true,
	}, nil
}

type PipelineConfig struct {
	DSI    dsihost.Config
	LTDC   ltdc.Config
	Layer  ltdc.LayerConfig
	Timing ltdc.Timing
}

func PipelineConfigForSurface(
	surface display.Surface,
	laneByteClockHz uint32,
	pixelClockHz uint32,
) (PipelineConfig, error) {
	if err := validateSurface(surface); err != nil {
		return PipelineConfig{}, err
	}

	timing := ltdc.Timing{
		Width:  uint16(surface.Width),
		Height: uint16(surface.Height),

		HSYNC: hSyncLen,
		HBP:   hBackPorch,
		HFP:   hFrontPorch,

		VSYNC: vSyncLen,
		VBP:   vBackPorch,
		VFP:   vFrontPorch,
	}

	dsi, err := DSIConfigForSurface(surface, laneByteClockHz, pixelClockHz)
	if err != nil {
		return PipelineConfig{}, err
	}

	pf, err := ltdcPixelFormat(surface.Format)
	if err != nil {
		return PipelineConfig{}, err
	}

	return PipelineConfig{
		DSI: dsi,

		Timing: timing,

		LTDC: ltdc.Config{
			Enable: true,
			Timing: timing,

			HSPolHigh: true,
			VSPolHigh: true,

			// LTDC side, not DSI side.
			DEPolHigh:          true,
			PixelClockInverted: true,

			Layers: [2]ltdc.LayerConfig{
				{
					Enabled:     true,
					Framebuffer: surface.Ptr,
					Width:       uint16(surface.Width),
					Height:      uint16(surface.Height),
					Pitch:       uint16(surface.Stride),
					Format:      pf,
					Alpha:       255,
				},
			},
		},
	}, nil
}
