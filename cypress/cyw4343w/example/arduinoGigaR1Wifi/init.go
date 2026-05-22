package main

import (
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/rcc"
)

//sigo:export preinit _preinit
func preinit() {
	// Disable the CM4.
	rcc.Rcc.Gcr.SetBootc2(false)

	// Configure the MPU.
	configureMPUForSDRAM()

	configureClocks()
	initSDRAM()
}

//sigo:export postinit _postinit
func postinit() {
	hal.SetDefaultFrequencies()
}
