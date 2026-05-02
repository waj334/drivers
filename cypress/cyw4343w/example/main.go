package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	_ "pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7"
	stm32h7x7 "pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/pin"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/sdio"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/timer"
	"pkg.si-go.dev/chip/arm/cortexm/platform/st/stm32h7x7/cm7/hal/uart"
	"pkg.si-go.dev/chip/arm/cortexm/runtime"
	"pkg.si-go.dev/drivers/cypress/cyw4343w"
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

//go:embed 4343WA1.bin
var firmware []byte

//go:embed 4343W1DX_CLM.bin
var clm []byte

//go:embed 4343W1DX_NVRAM.bin
var nvram []byte

var (
	WifiHost = cyw4343w.New[sdio.SDIO]()
	UART     = uart.UART1
)

func init() {
	// Prevent SysTick from driving timers.
	runtime.SysTickCanWake = false

	hal.ConfigureClocks()

	os.Stdout = &runtime.Semihosting

	// Configure UART
	err := UART.Configure(uart.Config{
		Enable: true,
		TX:     pin.PA9,
		RX:     pin.PB7,
		// FrameFormat:     uart.UsartFrame,
		Baud:            115_200,
		CharacterSize:   8,
		NumStopBits:     1,
		ReceiveEnabled:  true,
		TransmitEnabled: true,
	})

	if err != nil {
		fmt.Printf("Error configuring UART: %v\n", err)
		panic(err)
	}

	os.Stdout = UART

	err = TIM2.Configure(timer.Config{Enable: true})
	if err != nil {
		fmt.Printf("Error configuring TIM2: %v\n", err)
		panic(err)
	}
	stm32h7x7.IrqTim2.SetPriority(1)

	// Service the networking stack.
	go net.Poll()
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

	stm32h7x7.IrqSdmmc1.SetPriority(2)
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
		DMA: true,
	})

	if err != nil {
		fmt.Printf("Error configuring SDIO: %v\n", err)
		errorState()
		busyLoop()
	}

	err = WifiHost.Configure(cyw4343w.Config[sdio.SDIO]{
		Host:     SDIO1,
		Firmware: firmware,
		Nvram:    nvram,
		Clm:      clm,
	})

	if err != nil {
		fmt.Printf("Error configuring Wi-Fi host: %v\n", err)
		errorState()
		busyLoop()
	}

	// Initialize the card.
	err = WifiHost.InitializeCard()
	if err != nil {
		fmt.Printf("Error initializing Wi-Fi card: %v\n", err)
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
		fmt.Printf("Error reconfiguring SDIO: %v\n", err)
		errorState()
		busyLoop()
	}

	err = SDIO1.SetClockFrequency(10_000_000)
	if err != nil {
		fmt.Printf("Error setting SDIO clock frequency: %v\n", err)
		errorState()
		busyLoop()
	}

	// Initialize the Wi-Fi subsystem.
	err = WifiHost.Initialize()
	if err != nil {
		fmt.Printf("Error initializing Wi-Fi subsystem: %v\n", err)
		errorState()
		busyLoop()
	}

	// Start processing frames from the device.
	go func() {
		for {
			err := WifiHost.Poll()
			if err != nil {
				fmt.Printf("Error occurred while polling Wi-Fi subsystem: %v\n", err)
			}
		}
	}()

	// The CLM image needs to be loaded before any Wi-Fi/BLE functionality can be used.
	err = WifiHost.LoadClm()
	if err != nil {
		fmt.Printf("Error loading CLM: %v\n", err)
		errorState()
		busyLoop()
	}

	// Bring up the WLAN interface.
	err = WifiHost.Up()
	if err != nil {
		fmt.Printf("Error bringing up Wi-Fi interface: %v\n", err)
		errorState()
		busyLoop()
	}

	// Scan for Wi-Fi networks.
	networks, err := WifiHost.ScanWifiNetworks()
	if err != nil {
		fmt.Printf("Error scanning for Wi-Fi networks: %v\n", err)
		errorState()
		busyLoop()
	}

	for _, network := range networks {
		_, _ = os.Stdout.WriteString(network + "\n")
	}

	use(networks)

	// Join a WPA2 network.
	fmt.Printf("Joining WPA2 network...\n")
	err = WifiHost.JoinWPA2("waj334", "bigbluehooters")
	if err != nil {
		fmt.Printf("Error joining WPA2 network: %v\n", err)
		errorState()
		busyLoop()
	}
	fmt.Println("Joined!")

	// Print the MAC address.
	mac, err := WifiHost.MACAddress()
	if err != nil {
		fmt.Printf("Error getting MAC address: %v\n", err)
		errorState()
		busyLoop()
	} else {
		fmt.Printf("MAC: %02x:%02x:%02x:%02x:%02x:%02x\n",
			mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
	}

	// Register the WiFi driver as a net device to enable LWIP networking.
	// DHCP will start automatically.
	ni := net.RegisterNetDevice(WifiHost)

	// Wait for DHCP to assign an IP address.
	fmt.Println("Waiting for IP address...")
	for {
		ip := ni.IPAddress()
		if ip != ([4]byte{}) {
			fmt.Printf("IP: %d.%d.%d.%d\n", ip[0], ip[1], ip[2], ip[3])
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Dump google.com to Stdout.
	resp, err := http.Get("http://httpforever.com")
	if err != nil {
		fmt.Printf("GET error: %s\n", err)
		errorState()
		busyLoop()
	}

	_, err = io.Copy(os.Stdout, &resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("EOF")
		fmt.Printf("Copy error: %s\n", err)
		_ = resp.Body.Close()
		errorState()
		busyLoop()
	}

	_ = resp.Body.Close()

	goodState()
	busyLoop()
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
