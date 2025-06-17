package cyw4343w

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

type busFunctionType uint8

const (
	ioctlTypeGet = 0x00
	ioctlTypeSet = 0x02

	busFunction       busFunctionType = 0
	backplaneFunction busFunctionType = 1
	wlanFunction      busFunctionType = 2

	defaultTimeout = 100 * time.Millisecond
	htTimeout      = 2500 * time.Millisecond

	blockSize = 64
)

type Config[Host sdio.Host] struct {
	Host     Host
	Firmware []byte
	Nvram    []byte
	Clm      []byte
}

type Cyw4343w[Host sdio.Host] struct {
	host         Host
	receiveQueue queue

	requestIdCounter                  uint32
	bsscIndex                         uint32 // Leaving at 0 for now.
	backplaneWindowCurrentBaseAddress uint32
	chipId                            uint32
	txSeq                             uint8
	txMax                             uint8

	firmware []byte
	nvram    []byte
	clm      []byte

	mutex      sync.Mutex
	iovarMutex sync.Mutex
	queueMutex sync.Mutex
}

func New[Host sdio.Host]() *Cyw4343w[Host] {
	// Create an instance of the driver.
	c := &Cyw4343w[Host]{
		requestIdCounter: 1,
		txSeq:            1,
	}

	// Initialize the queue.
	c.receiveQueue = queue{
		mutex: &c.queueMutex,
	}

	return c
}

func (c *Cyw4343w[SDIO]) Configure(config Config[SDIO]) error {
	c.mutex.Lock()

	if len(config.Firmware) == 0 {
		c.mutex.Unlock()
		return errInvalidFirmwareImage
	}

	if len(config.Nvram) == 0 {
		c.mutex.Unlock()
		return errInvalidNvramImage
	}

	if len(config.Clm) == 0 {
		c.mutex.Unlock()
		return errInvalidClmImage
	}

	c.host = config.Host
	c.firmware = config.Firmware
	c.nvram = config.Nvram
	c.clm = config.Clm

	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) Initialize() error {
	c.mutex.Lock()

	// Initialize the backplane.
	err := c.initBackplane()
	if err != nil {
		c.mutex.Unlock()
		return err
	}

	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) InitializeCard() error {
	c.mutex.Lock()

	var resp sdio.Response
	var err error
	ready := false

	// Send CMD0 to reset the card
	_, err = c.host.SendCommand(sdio.Command{
		Class: sdio.CMD0,
	})

	if err != nil {
		c.mutex.Unlock()
		return err
	}

	for retry := 0; retry < 1000; retry++ {
		// Send CMD5 to get the ready status of the card.
		resp, err = c.host.SendCommand(sdio.Command{
			Class: sdio.CMD5,
			Argument: sdio.CMD5Args{
				VoltageBits: 0x00FF8000,
			}.Value(),
		})

		if err != nil {
			time.Sleep(time.Millisecond * 5)
			continue
		} else {
			r4 := sdio.R4(resp[0])
			if !r4.Ready() {
				// The card is not ready. Try again.
				time.Sleep(time.Millisecond * 5)
				continue
			}
			ready = true
		}
		break
	}

	if !ready {
		c.mutex.Unlock()
		return sdio.ErrNotReady
	}

	// Send CMD3 to get the address of the card.
	resp, err = c.host.SendCommand(sdio.Command{
		Class: sdio.CMD3,
	})

	if err != nil {
		c.mutex.Unlock()
		return err
	}

	r6 := sdio.R6(resp[0])

	// Send CMD7 with the returned RCA to select the card.
	_, err = c.host.SendCommand(sdio.Command{
		Class: sdio.CMD7,
		Argument: sdio.CMD7Args{
			RCA: r6.RCA(),
		}.Value(),
	})

	if err != nil {
		c.mutex.Unlock()
		return err
	}

	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) initBackplane() error {
	// Enable the backplane.
	deadline := time.Now().Add(defaultTimeout)
	for {
		err := c.writeRegisterValue(busFunction, sdiodCccrIoen, 1, sdioFuncEnable1)
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond)

		var value uint32
		value, err = c.readRegisterValue(busFunction, sdiodCccrIoen, 1)
		if err != nil {
			return err
		}

		if value == sdioFuncEnable1 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}
	}

	// Wait until the backplane is ready.
	deadline = time.Now().Add(defaultTimeout)
	for {
		err := c.writeRegisterValue(busFunction, sdiodCccrBlksize0, 1, 64)
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond)

		value, err := c.readRegisterValue(busFunction, sdiodCccrBlksize0, 1)
		if err != nil {
			return err
		}

		if value == 64 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}
	}

	// Set all block register sizes.
	err := c.writeRegisterValue(busFunction, sdiodCccrBlksize0, 1, blockSize)
	if err != nil {
		return err
	}

	err = c.writeRegisterValue(busFunction, sdiodCccrF1blksize0, 1, blockSize)
	if err != nil {
		return err
	}

	err = c.writeRegisterValue(busFunction, sdiodCccrF2blksize0, 1, blockSize)
	if err != nil {
		return err
	}

	err = c.writeRegisterValue(busFunction, sdiodCccrF2blksize1, 1, 0)
	if err != nil {
		return err
	}

	// Enable interrupts.
	err = c.writeRegisterValue(busFunction, sdiodCccrInten, 1, 0x1|0x2|0x4)
	if err != nil {
		return err
	}

	// Wait for the backplane to be ready.
	deadline = time.Now().Add(defaultTimeout)
	for {
		value6, err := c.readRegisterValue(busFunction, sdiodCccrIordy, 1)
		if err != nil {
			return err
		}

		// Function 1 IO Ready.
		if value6&0x02 != 0 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}

		time.Sleep(time.Millisecond)
	}

	// Set the ALP.
	err = c.writeRegisterValue(backplaneFunction, sdioChipClockCsr, 1, 0x08)
	if err != nil {
		return err
	}

	for {
		value3, err := c.readRegisterValue(backplaneFunction, sdioChipClockCsr, 1)
		if err != nil {
			return err
		}

		if value3&0x40 != 0 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}
	}

	// Clear ALP request.
	err = c.writeRegisterValue(backplaneFunction, sdioChipClockCsr, 1, 0)
	if err != nil {
		return err
	}

	// Disable extra SDIO pull-ups.
	err = c.writeRegisterValue(backplaneFunction, sdioPullUp, 1, 0)
	if err != nil {
		return err
	}

	// Enable F1 and F2
	err = c.writeRegisterValue(busFunction, sdiodCccrIoen, 1, sdioFuncEnable1|sdioFuncEnable2)
	if err != nil {
		return err
	}

	// TODO: Set up host-wake signals.

	// Enable F2 interrupt only.
	err = c.writeRegisterValue(backplaneFunction, sdiodCccrInten, 1, 0x01|0x04)
	if err != nil {
		return err
	}

	// Check if the chip supports CHIPID read from SDIO core and bootloader handshake.
	value7, err := c.readRegisterValue(busFunction, sdiodCccrBrcmCardcap, 1)
	if value7&0x40 != 0 {
		var addrLow, addrMid, addrHigh, devctl uint32
		var err error

		devctl, err = c.readRegisterValue(backplaneFunction, sbsdioDeviceCtl, 1)
		if err != nil {
			return err
		}

		err = c.writeRegisterValue(backplaneFunction, sbsdioDeviceCtl, 1, devctl|sbsdioDevctlAddrRst)
		if err != nil {
			return err
		}

		addrLow, err = c.readRegisterValue(backplaneFunction, sbsdioFunc1Sbaddrlow, 1)
		if err != nil {
			return err
		}

		addrMid, err = c.readRegisterValue(backplaneFunction, sbsdioFunc1Sbaddrmid, 1)
		if err != nil {
			return err
		}

		addrHigh, err = c.readRegisterValue(backplaneFunction, sbsdioFunc1Sbaddrhigh, 1)
		if err != nil {
			return err
		}

		regAddr := ((addrLow << 8) | (addrMid << 16) | (addrHigh << 24)) + sdioCoreChipidReg
		err = c.writeRegisterValue(backplaneFunction, sbsdioDeviceCtl, 1, devctl)
		if err != nil {
			return err
		}

		// Read the chip id.
		c.chipId, err = c.readBackplaneValue(regAddr, 2)
		if err != nil {
			return err
		}
	} else {
		c.chipId, err = c.readBackplaneValue(chipcommonBaseAddress, 2)
		if err != nil {
			return err
		}
	}

	err = c.downloadFirmware()
	if err != nil {
		return err
	}

	// Wait for F2 to be ready.
	deadline = time.Now().Add(defaultTimeout)
	for {
		value8, err := c.readRegisterValue(busFunction, sdiodCccrIordy, 1)
		if err != nil {
			return err
		}

		if value8&0x04 != 0 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}
	}

	err = c.enableSaveRestore()
	if err != nil {
		return err
	}

	// Poll for initial credits.
	err = c.poll()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) LoadClm() error {
	type clmHeaderType struct {
		flag   uint16
		typ    uint16
		length uint32
		crc    uint32
	}

	if len(c.clm) == 0 {
		return errInvalidClmImage
	}

	clm := c.clm
	const maxLoadLen = 512
	const clmHeaderLength = int(unsafe.Sizeof(clmHeaderType{}))
	var buffer [maxLoadLen + clmHeaderLength]byte

	offset := 0
	header := (*clmHeaderType)(unsafe.Pointer(&buffer[0]))
	*header = clmHeaderType{
		flag: 1<<12 | 2,
		typ:  2,
		crc:  0,
	}

	for offset < len(clm) {
		n := min(maxLoadLen, len(clm)-offset)
		header.length = uint32(n)

		if offset+n >= len(clm) {
			// Set the last chunk flag.
			header.flag |= 4
		}

		// Copy the data to the buffer after the header.
		copy(buffer[clmHeaderLength:], clm[:n])
		clm = clm[n:]

		// Send the buffer to the chip.
		_, err := c.SetIovar(IovarStrClmload, buffer[:clmHeaderLength+n])
		if err != nil {
			return err
		}
		offset += n

		// Reset flags.
		header.flag = 1 << 12
	}

	return nil
}

func (c *Cyw4343w[SDIO]) downloadFirmware() error {
	ramStartAddress := atcmRamBaseAddress(c.chipId)
	if ramStartAddress != 0 {
		// TODO: Reset the WLAN ARM core.
	} else {
		err := c.disableDeviceCore(wlanArmCore, wlanCoreFlagNone)
		if err != nil {
			return err
		}

		err = c.disableDeviceCore(socramCore, wlanCoreFlagNone)
		if err != nil {
			return err
		}

		err = c.resetDeviceCore(socramCore, wlanCoreFlagNone)
		if err != nil {
			return err
		}

		err = c.chipSpecificSocsramInit()
		if err != nil {
			return err
		}
	}

	firmware := c.firmware
	nvram := c.nvram
	ramAddr := atcmRamBaseAddress(c.chipId)
	ramSize := chipRamSize(c.chipId)

	err := c.writeBackplaneBytes(ramAddr, firmware)
	if err != nil {
		return err
	}

	nvramLen := roundUp(uint32(len(nvram)), 4)
	err = c.writeBackplaneBytes(ramAddr+ramSize-4-nvramLen, nvram)
	if err != nil {
		return err
	}

	nvramLenWords := nvramLen / 4
	nvramLenMagic := ((^nvramLenWords) << 16) | nvramLenWords
	err = c.writeBackplaneValue(ramAddr+ramSize-4, 4, nvramLenMagic)
	if err != nil {
		return err
	}

	if ramStartAddress != 0 {
		// TODO: Reset core with bits.
	} else {
		err = c.resetDeviceCore(wlanArmCore, wlanCoreFlagNone)
		if err != nil {
			return err
		}

		isUp, err := c.deviceCoreIsUp(wlanArmCore)
		if err != nil {
			return err
		}

		if !isUp {
			return errCoreIsNotUp
		}
	}

	// Wait until high throughput clock is ready.
	deadline := time.Now().Add(htTimeout)
	for {
		csr, err := c.readRegisterValue(backplaneFunction, sdioChipClockCsr, 1)
		if err != nil {
			return err
		}

		if csr&sbsdioHtAvail != 0 {
			break
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}

		time.Sleep(time.Millisecond)
	}

	// set up the interrupt mask and enable interrupts.
	base := sdiodCoreBaseAddress(c.chipId)
	err = c.writeBackplaneValue(base+0x24, 4, 0x000000F0)
	if err != nil {
		return err
	}

	err = c.writeBackplaneValue(base+0x34, 1, 0x01|0x02)
	if err != nil {
		return err
	}

	err = c.writeRegisterValue(backplaneFunction, sdioFunction2Watermark, 1, 8)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) writeBackplaneValue(address uint32, length uint8, value uint32) error {
	err := c.setBackplaneWindow(address)
	if err != nil {
		return err
	}

	address &= sbOftAddrMask
	if length == 4 {
		address |= sbAccess24BFlag
	}

	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], value)
	err = c.write(DataTransfer{
		Data:     data[:min(4, length)],
		Address:  address,
		Function: uint32(backplaneFunction),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Cyw4343w[SDIO]) readBackplaneValue(address uint32, length uint8) (uint32, error) {
	err := c.setBackplaneWindow(address)
	if err != nil {
		return 0, err
	}

	address &= sbOftAddrMask
	if length == 4 {
		address |= sbAccess24BFlag
	}

	var data [4]byte
	err = c.read(DataTransfer{
		Data:     data[:min(4, length)],
		Address:  address,
		Function: uint32(backplaneFunction),
	})
	if err != nil {
		return 0, err
	}

	var mask uint32
	switch length {
	case 1:
		mask = 0x0000_00FF
	case 2:
		mask = 0x0000_FFFF
	case 3:
		mask = 0x00FF_FFFF
	default:
		mask = 0xFFFF_FFFF
	}

	return binary.LittleEndian.Uint32(data[:]) & mask, nil
}

func (c *Cyw4343w[SDIO]) writeBackplaneBytes(address uint32, data []byte) error {
	remaining := uint32(len(data))
	offset := uint32(0)
	for remaining > 0 {
		// Determine the transfer size for this chunk.
		transferSize := uint32(maxBackplaneTransferSize)
		if remaining < transferSize {
			transferSize = remaining
		}

		// Make sure we don't cross the backplane window boundary.
		windowOffset := address & backplaneAddressMask
		if windowOffset+transferSize >= backplaneWindowSize {
			transferSize = backplaneWindowSize - windowOffset
		}

		// Set backplane window.
		err := c.setBackplaneWindow(address)
		if err != nil {
			return err
		}

		// DataTransfer address within the window.
		transAddr := address & backplaneAddressMask
		chunk := data[offset : offset+transferSize]
		err = c.write(DataTransfer{
			Data:     chunk,
			Address:  transAddr,
			Function: uint32(backplaneFunction),
		})
		if err != nil {
			return err
		}

		address += transferSize
		offset += transferSize
		remaining -= transferSize
	}

	// Reset the backplane window address.
	err := c.setBackplaneWindow(0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) readBackplaneBytes(address uint32, data []byte) error {
	remaining := uint32(len(data))
	offset := uint32(0)
	for remaining > 0 {
		// Determine the transfer size for this chunk.
		transferSize := uint32(maxBackplaneTransferSize)
		if remaining < transferSize {
			transferSize = remaining
		}

		// Make sure we don't cross the backplane window boundary.
		windowOffset := address & backplaneAddressMask
		if windowOffset+transferSize > backplaneAddressMask {
			transferSize = backplaneWindowSize - windowOffset
		}

		// Set backplane window.
		err := c.setBackplaneWindow(address)
		if err != nil {
			return err
		}

		// DataTransfer address within the window.
		transAddr := address & backplaneAddressMask

		// For reads, we must allocate buffer with padding.
		readSize := transferSize + backplaneReadPaddSize
		readBuf := make([]byte, readSize)
		err = c.read(DataTransfer{
			Data:     readBuf,
			Address:  transAddr,
			Function: uint32(backplaneFunction),
		})
		if err != nil {
			return err
		}

		// Copy out payload, skipping padding
		copy(data[offset:offset+transferSize], readBuf[backplaneReadPaddSize:])

		address += transferSize
		offset += transferSize
		remaining -= transferSize
	}

	// Reset the backplane window address.
	err := c.setBackplaneWindow(0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) setBackplaneWindow(addr uint32) error {
	const upperMask32 = uint32(0xFF00_0000)
	const upperMiddleMask32 = uint32(0x00FF_0000)
	const lowerMiddleMask32 = uint32(0x0000_FF00)

	base := addr & ^uint32(backplaneAddressMask)
	if base == c.backplaneWindowCurrentBaseAddress {
		return nil
	}

	if (base & upperMask32) != (c.backplaneWindowCurrentBaseAddress & upperMask32) {
		// Write register value.
		err := c.writeRegisterValue(backplaneFunction, sdioBackplaneAddressHigh, 1, base>>24)
		if err != nil {
			return err
		}

		// Clear old value.
		c.backplaneWindowCurrentBaseAddress &= ^upperMask32

		// Set new value.
		c.backplaneWindowCurrentBaseAddress |= base & upperMask32
	}

	if (base & upperMiddleMask32) != (c.backplaneWindowCurrentBaseAddress & upperMiddleMask32) {
		// Write register value.
		err := c.writeRegisterValue(backplaneFunction, sdioBackplaneAddressMid, 1, base>>16)
		if err != nil {
			return err
		}

		// Clear old value.
		c.backplaneWindowCurrentBaseAddress &= ^upperMiddleMask32

		// Set new value.
		c.backplaneWindowCurrentBaseAddress |= base & upperMiddleMask32
	}

	if (base & lowerMiddleMask32) != (c.backplaneWindowCurrentBaseAddress & lowerMiddleMask32) {
		// Write register value.
		err := c.writeRegisterValue(backplaneFunction, sdioBackplaneAddressLow, 1, base>>8)
		if err != nil {
			return err
		}

		// Clear old value.
		c.backplaneWindowCurrentBaseAddress &= ^lowerMiddleMask32

		// Set new value.
		c.backplaneWindowCurrentBaseAddress |= base & lowerMiddleMask32
	}

	return nil
}

func (c *Cyw4343w[SDIO]) writeRegisterValue(function busFunctionType, address uint32, length uint8, value uint32) error {
	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], value)
	return c.write(DataTransfer{
		Data:     data[:min(4, length)],
		Address:  address,
		Function: uint32(function),
	})
}

func (c *Cyw4343w[SDIO]) readRegisterValue(function busFunctionType, address uint32, length uint8) (uint32, error) {
	var data [4]byte
	err := c.read(DataTransfer{
		Data:     data[:min(4, length)],
		Address:  address,
		Function: uint32(function),
	})

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data[:]), nil
}

func (c *Cyw4343w[SDIO]) requestId() uint32 {
	return atomic.AddUint32(&c.requestIdCounter, 1)
}

func (c *Cyw4343w[SDIO]) chipSpecificSocsramInit() error {
	switch c.chipId {
	case 43430, 43439:
		const sramBase = 0x18000000 + 0x4000
		err := c.writeBackplaneValue(sramBase+0x10, 4, 0x3)
		if err != nil {
			return err
		}

		err = c.writeBackplaneValue(sramBase+0x44, 4, 0)
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func (c *Cyw4343w[SDIO]) enableSaveRestore() error {
	ok, err := c.isFwSrCapable()
	if err != nil {
		return err
	}

	if ok {
		// Configure the WakeupCtrl register to set HtAvail request bit in chipClockCSR register after the sdiod core
		// is powered on.
		data, err := c.readRegisterValue(backplaneFunction, sdioWakeupCtrl, 1)
		if err != nil {
			return err
		}
		data |= sbsdioWctrlWlWakeTillHtAvail
		err = c.writeRegisterValue(backplaneFunction, sdioWakeupCtrl, 1, data)
		if err != nil {
			return err
		}
	}

	// Set brcmCardCapability to noCmdDecode mode. It makes sdiod_aos able to wake up the host for any activity of cmd
	// line, even though the module won't decode or respond.
	err = c.writeRegisterValue(busFunction, sdiodCccrBrcmCardcap, 1, sdiodCccrBrcmCardcapCmdNodec)
	if err != nil {
		return err
	}

	err = c.writeRegisterValue(backplaneFunction, sdioChipClockCsr, 1, sbsdioForceHt)
	if err != nil {
		return err
	}

	// Enable KeepSdioOn (KSO) bit for normal operation.
	data, err := c.readRegisterValue(backplaneFunction, sdioSleepCsr, 1)
	if err != nil {
		return err
	}

	if data&sbsdioSlpcsrKeepWlKso == 0 {
		data |= sbsdioSlpcsrKeepWlKso
		err = c.writeRegisterValue(backplaneFunction, sdioSleepCsr, 1, data)
	}

	// Put SPI interface block to sleep.
	err = c.writeRegisterValue(backplaneFunction, sdioPullUp, 1, 0xF)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) isFwSrCapable() (bool, error) {
	srCtrl, err := c.readBackplaneValue(chipcommonBaseAddress+0x508, 4)
	if err != nil {
		return false, err
	}
	return srCtrl != 0, nil
}
