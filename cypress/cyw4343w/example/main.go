package main

import (
	"pkg.si-go.dev/drivers/cypress/cyw4343w"
	"time"

	_ "pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/pin"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/sdio"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/timer"
	"pkg.si-go.dev/chip/arm/cortexm/runtime"
)

const (
	LEDR  = pin.PI12
	LEDB  = pin.PE3
	LEDG  = pin.PJ13
	TIM2  = timer.TIM2
	SDIO1 = sdio.SDIO1

	BLEON = pin.PA10

	WIFION       = pin.PB10
	WIFIHOSTWAKE = pin.PI8

	SDIO_D1  = pin.PC8
	SDIO_D2  = pin.PC9
	SDIO_D3  = pin.PC10
	SDIO_D4  = pin.PC11
	SDIO_CLK = pin.PC12
	SDIO_CMD = pin.PD2

	timescale = uint64(time.Microsecond)
)

var (
	WifiHost = cyw4343w.New[sdio.SDIO]()
)

//sigo:export wake runtime.wake
func wake(t uint64)

func alarm(t uint64) {
	wake(t)
}

//sigo:export nanotime runtime.nanotime
func nanotime() uint64 {
	// The timer resolution is 1uS per tick.
	return TIM2.Tick() * timescale
}

//sigo:export addsleep runtime.addsleep
func addsleep(deadline uint64) {
	TIM2.SetAlarm(deadline/timescale, alarm)
}

func init() {
	// Prevent SysTick from driving timers.
	runtime.SysTickCanWake = false

	hal.ConfigureClocks()

	err := TIM2.Configure(timer.Config{Enable: true})
	if err != nil {
		panic(err)
	}
}

func main() {
	LEDR.SetMode(pin.Output)
	LEDB.SetMode(pin.Output)
	LEDG.SetMode(pin.Output)

	LEDR.High()
	LEDB.Low()
	LEDG.High()

	SDIO_D1.SetSpeedMode(pin.VeryHighSpeed)
	SDIO_D2.SetSpeedMode(pin.VeryHighSpeed)
	SDIO_D3.SetSpeedMode(pin.VeryHighSpeed)
	SDIO_D4.SetSpeedMode(pin.VeryHighSpeed)
	SDIO_CLK.SetSpeedMode(pin.VeryHighSpeed)
	SDIO_CMD.SetSpeedMode(pin.VeryHighSpeed)

	SDIO_D1.SetPullMode(pin.PullUp)
	SDIO_D2.SetPullMode(pin.PullUp)
	SDIO_D3.SetPullMode(pin.PullUp)
	SDIO_D4.SetPullMode(pin.PullUp)
	SDIO_CLK.SetPullMode(pin.PullUp)
	SDIO_CMD.SetPullMode(pin.PullUp)

	WIFIHOSTWAKE.SetMode(pin.Input)
	WIFION.SetMode(pin.Output)
	WIFION.High()

	BLEON.SetMode(pin.Output)
	BLEON.High()

	time.Sleep(time.Millisecond * 250)

	err := SDIO1.Configure(sdio.Config{
		Enable: true,
		CK:     SDIO_CLK,
		Dn: [8]pin.Pin{
			SDIO_D1,
			SDIO_D2,
			SDIO_D3,
			SDIO_D4,
		},
		CMD: SDIO_CMD,
		DMA: false,
	})

	if err != nil {
		errorState()
		busyLoop()
	}

	err = WifiHost.Configure(cyw4343w.Config[sdio.SDIO]{Host: SDIO1})
	if err != nil {
		errorState()
		busyLoop()
	}

	// Initialize the card.
	err = WifiHost.InitializeCard()
	if err != nil {
		errorState()
		busyLoop()
	}

	// Reconfigure the SDIO host interface.
	err = SDIO1.Reconfigure(sdio.SecondaryConfig{
		NegEdge:             false,
		BusWidth:            sdio.Width4Bit,
		BusSpeed:            sdio.Hs,
		PowerSave:           false,
		HardwareFlowControl: false,
	})

	if err != nil {
		errorState()
		busyLoop()
	}

	err = SDIO1.SetClockFrequency(5_000_000)
	if err != nil {
		errorState()
		busyLoop()
	}

	// Initialize the Wi-Fi subsystem.
	err = WifiHost.Initialize()
	if err != nil {
		errorState()
		busyLoop()
	}

	goodState()

	// Start processing frames from the device.
	go func() {
		for {
			err := WifiHost.Poll()
			if err != nil {
				//errorState()
				//return
			}
		}
	}()

	for {
		response, err := WifiHost.Iovar(cyw4343w.IovarStrVersion, 0)
		if err != nil {
			errorState()
		}
		use(response)
	}
}

func errorState() {
	LEDR.Low()
	LEDB.High()
	LEDG.High()
}

func goodState() {
	LEDR.High()
	LEDB.High()
	LEDG.Low()
}

func busyLoop() {
	select {}
}

func use(any) {}
