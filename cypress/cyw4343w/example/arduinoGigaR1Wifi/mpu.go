package main

import (
	"pkg.si-go.dev/chip/arm/cortexm/reg/mpu"
	cortexm "pkg.si-go.dev/chip/arm/cortexm/runtime"
)

const (
	sdramBase   = 0xC0000000
	sdramRegion = uint8(0)
	sdramSize   = uint8(22)    // 8MB
	sdramAp     = uint8(0b011) // FullAccess
	sdramTex    = uint8(0b001) // Normal memory, non-cacheable
)

func configureMPUForSDRAM() {
	// Disable MPU first.
	mpu.Mpu.Ctrl.SetEnable(false)
	cortexm.DSB()
	cortexm.ISB()

	// Region 0: SDRAM 8MB at 0xC0000000.
	mpu.Mpu.Rnr.SetRegion(sdramRegion)

	// Base address for selected region.
	mpu.Mpu.Rbar.SetAddr(sdramBase >> 5)

	// Region attributes.
	mpu.Mpu.Rasr.SetXn(false) // allow execution
	mpu.Mpu.Rasr.SetAp(sdramAp)
	mpu.Mpu.Rasr.SetTex(sdramTex)
	mpu.Mpu.Rasr.SetS(false) // non-shareable
	mpu.Mpu.Rasr.SetC(false) // non-cacheable
	mpu.Mpu.Rasr.SetB(false) // non-bufferable
	mpu.Mpu.Rasr.SetSrd(0)   // no subregions disabled
	mpu.Mpu.Rasr.SetSize(sdramSize)
	mpu.Mpu.Rasr.SetEnable(true)

	// Enable MPU with PRIVDEFENA so everything else uses the default memory map.
	mpu.Mpu.Ctrl.SetPrivdefena(true)
	mpu.Mpu.Ctrl.SetEnable(true)

	cortexm.DSB()
	cortexm.ISB()
}
