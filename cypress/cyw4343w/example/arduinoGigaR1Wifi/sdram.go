package main

import (
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/pin"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/sdram"
)

// SDRAM bring-up for the Arduino GIGA R1.
//
// AS4C4M16SA-7 or M12L64164A (4M × 16, 8 MB total) on FMC Bank 1, mapped at 0xC0000000.
func initSDRAM() {
	_ = sdram.Configure(sdram.Config{
		Bank: [2]sdram.BankConfig{
			{
				Enable:          true,
				BusWidth:        sdram.BusWidth16,
				RowBits:         sdram.RowBits12Bit,
				ColumnBits:      sdram.ColumnBits8Bit,
				InternalBanks:   sdram.BankCount4,
				BurstRead:       true,
				WriteProtection: false,
				ReadPipeDelay:   sdram.NoDelay,
				CASLatency:      sdram.CAS3,
				ClockPeriod:     sdram.Period2Cycles,
				LoadMode:        0x230,
				Timing: sdram.TimingConfig{
					TRCD: 3,
					TRP:  3,
					TWR:  3,
					TRC:  7,
					TRAS: 5,
					TXSR: 7,
					TMRD: 2,
				},
			},
		},
		CLK: pin.PG8,
		CKE: pin.PH2,
		CS:  pin.PH3,
		RAS: pin.PF11,
		CAS: pin.PG15,
		WE:  pin.PH5,
		ADDR: [13]pin.Pin{
			pin.PF0, pin.PF1, pin.PF2, pin.PF3,
			pin.PF4, pin.PF5, pin.PF12, pin.PF13,
			pin.PF14, pin.PF15, pin.PG0, pin.PG1,
			pin.PG2,
		},
		DATA: [32]pin.Pin{
			pin.PD14, pin.PD15, pin.PD0, pin.PD1,
			pin.PE7, pin.PE8, pin.PE9, pin.PE10,
			pin.PE11, pin.PE12, pin.PE13, pin.PE14,
			pin.PE15, pin.PD8, pin.PD9, pin.PD10,
		},
		BA: [2]pin.Pin{
			pin.PG4, pin.PG5,
		},
		NBL: [4]pin.Pin{
			pin.PE0, pin.PE1,
		},
		RT:                 0x0000C0C,
		PowerUpDelayCycles: 200_000,
	})
}
