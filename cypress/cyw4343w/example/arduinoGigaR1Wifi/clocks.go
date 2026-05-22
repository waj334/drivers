package main

import (
	"volatile"

	cortexm "pkg.si-go.dev/chip/arm/cortexm/runtime"

	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/flash"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/pwr"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/rcc"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/reg/syscfg"
)

// Captured RCC register values from CubeMX for the target configuration:
//
//	PLL1: HSI/4 ×60 /2 = 480 MHz CPU, /8 = 120 MHz Q, /2 = 480 MHz R
//	PLL2: HSI/32 ×300 /2 = 300 MHz P, /2 = 300 MHz Q, /3 = 200 MHz R (FMC/SDMMC)
//	PLL3: configured but disabled
//	HCLK = 240 MHz, PCLKs = 120 MHz
//	FMC source = PLL2R (200 MHz)
//	SDMMC source = PLL1Q (120 MHz)  ← change if you want both on PLL2R
const (
	rccCrPllsOn     = uint32(0x05001000) // PLL1ON + PLL2ON + HSI48ON
	rccPllckselrVal = uint32(0x02020040) // HSI source, Divm1=4, Divm2=32, Divm3=16
	rccPllcfgrVal   = uint32(0x01FF005D) // Output enables + input ranges
	rccPll1divrVal  = uint32(0x0107023B) // Divn1=60, Divp1=2, Divq1=8, Divr1=2
	rccPll2divrVal  = uint32(0x0201032B) // Divn2=300, Divp2=2, Divq2=2, Divr2=3 (PLL2R = 200 MHz)
	rccPll3divrVal  = uint32(0x01010280) // Divn3=129, Divp3=2, Divq3=2, Divr3=2
	rccD1cfgrVal    = uint32(0x00000048) // D1CPRE=/1, HPRE=/2, D1PPRE=/2
	rccD2cfgrVal    = uint32(0x00000440) // D2PPRE1=/2, D2PPRE2=/2
	rccD3cfgrVal    = uint32(0x00000040) // D3PPRE=/2
	rccCfgrSwPll1   = uint32(0x00000003) // SW = PLL1
	rccD1ccipirVal  = uint32(0x00010002) // FMCSRC=PLL2R, SDMMCSRC=PLL1Q
	rccD2ccip1rVal  = uint32(0x00000000) // SPI/SAI defaults (PLL1Q for SPI123)
	rccD2ccip2rVal  = uint32(0x00000000) // USART uses PCLK
	rccD3ccipirVal  = uint32(0x00000000)
)

// ConfigureClocks programs the H7x7 clock tree to the configuration captured
// from CubeMX. Writes only constants and hardware registers — safe to call
// from preinit before .data/.bss are initialized.
func configureClocks() {
	state := cortexm.DisableInterrupts()

	pwr.Pwr.Cr3.SetSden(false)
	for !pwr.Pwr.Csr1.GetActvosrdy() {
	}

	rcc.Rcc.Apb4enr.SetSyscfgen(true)
	_ = rcc.Rcc.Apb4enr.GetSyscfgen()
	cortexm.DSB()

	pwr.Pwr.D3cr.SetVos(pwr.RegisterD3crFieldVosEnumScale1)
	for !pwr.Pwr.D3cr.GetVosrdy() {
	}

	syscfg.Syscfg.Pwrcr.SetOden(true)
	for !pwr.Pwr.D3cr.GetVosrdy() {
	}

	// Configure all PLLs (PLLs must be off — they are at reset).
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Pllckselr), rccPllckselrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Pllcfgr), rccPllcfgrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Pll1divr), rccPll1divrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Pll2divr), rccPll2divrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Pll3divr), rccPll3divrVal)

	// Bus prescalers — set before SYSCLK switches to PLL1 to avoid overshoot.
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D1cfgr), rccD1cfgrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D2cfgr), rccD2cfgrVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D3cfgr), rccD3cfgrVal)

	// Flash wait states for 480 MHz at VOS0.
	flash.Flash.Bank[0].Acr.SetLatency(4)
	for flash.Flash.Bank[0].Acr.GetLatency() != 4 {
	}
	flash.Flash.Bank[0].Acr.SetWrhighfreq(3)
	for flash.Flash.Bank[0].Acr.GetWrhighfreq() != 3 {
	}

	// Enable PLL1 and PLL2. PLL3 stays off.
	cr := volatile.LoadUint32((*uint32)(&rcc.Rcc.Cr))
	cr |= rccCrPllsOn
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Cr), cr)

	// Wait for PLL1, PLL2, and HSI48 to be ready.
	for volatile.LoadUint32((*uint32)(&rcc.Rcc.Cr))&(1<<25) == 0 { // PLL1RDY
	}
	for volatile.LoadUint32((*uint32)(&rcc.Rcc.Cr))&(1<<27) == 0 { // PLL2RDY
	}
	for volatile.LoadUint32((*uint32)(&rcc.Rcc.Cr))&(1<<13) == 0 { // RC48RDY (HSI48)
	}

	// Switch SYSCLK to PLL1.
	cfgr := volatile.LoadUint32((*uint32)(&rcc.Rcc.Cfgr))
	cfgr = (cfgr &^ 0x7) | rccCfgrSwPll1
	volatile.StoreUint32((*uint32)(&rcc.Rcc.Cfgr), cfgr)
	for (volatile.LoadUint32((*uint32)(&rcc.Rcc.Cfgr))>>3)&0x7 != rccCfgrSwPll1 {
	}

	// Peripheral kernel clock muxes.
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D1ccipr), rccD1ccipirVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D2ccip1r), rccD2ccip1rVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D2ccip2r), rccD2ccip2rVal)
	volatile.StoreUint32((*uint32)(&rcc.Rcc.D3ccipr), rccD3ccipirVal)

	// Peripheral clock enables for this application:
	//   GPIO A-K (all banks for FMC SDRAM + SDIO + UART pins)
	//   FMC (SDRAM)
	//   SDMMC1 (WiFi)
	//   USART1 (console)
	rcc.Rcc.Ahb4enr.SetGpioaen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioaen()
	rcc.Rcc.Ahb4enr.SetGpioben(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioben()
	rcc.Rcc.Ahb4enr.SetGpiocen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpiocen()
	rcc.Rcc.Ahb4enr.SetGpioden(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioden()
	rcc.Rcc.Ahb4enr.SetGpioeen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioeen()
	rcc.Rcc.Ahb4enr.SetGpiofen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpiofen()
	rcc.Rcc.Ahb4enr.SetGpiogen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpiogen()
	rcc.Rcc.Ahb4enr.SetGpiohen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpiohen()
	rcc.Rcc.Ahb4enr.SetGpioien(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioien()
	rcc.Rcc.Ahb4enr.SetGpiojen(true)
	_ = rcc.Rcc.Ahb4enr.GetGpiojen()
	rcc.Rcc.Ahb4enr.SetGpioken(true)
	_ = rcc.Rcc.Ahb4enr.GetGpioken()

	rcc.Rcc.Ahb3enr.SetFmcen(true)
	_ = rcc.Rcc.Ahb3enr.GetFmcen()
	rcc.Rcc.Ahb3enr.SetSdmmc1en(true)
	_ = rcc.Rcc.Ahb3enr.GetSdmmc1en()
	rcc.Rcc.Apb2enr.SetUsart1en(true)
	_ = rcc.Rcc.Apb2enr.GetUsart1en()
	rcc.Rcc.Ahb2enr.SetRngen(true)
	_ = rcc.Rcc.Ahb2enr.GetRngen()

	cortexm.DSB()
	cortexm.EnableInterrupts(state)
}
