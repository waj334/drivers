package cyw4343w

import (
	"encoding/binary"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal"
)

func ioTransfer(data []byte) DataTransfer {
	return DataTransfer{
		Data:     data,
		Address:  0,
		Function: uint32(wlanFunction),
	}
}

func (c *Cyw4343w[SDIO]) Poll() error {
	c.mutex.Lock()
	err := c.poll()
	c.mutex.Unlock()
	return err
}

func (c *Cyw4343w[SDIO]) poll() error {
	// Check if the interrupt indicated there is a packet to read.
	available, err := c.packetAvailableToRead()
	if err != nil {
		return err
	}

	if available {
		err = c.receiveOnePacket()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cyw4343w[SDIO]) readFrame() ([]byte, error) {
	// TODO: Check that the WLAN backplane is up before continuing...

	var hwTag [2]uint16
	bhwTag := unsafe.Slice((*byte)(unsafe.Pointer(&hwTag[0])), 4)

	// Read the frame header and verify validity.
	err := c.read(ioTransfer(bhwTag[:4]))
	if err != nil {
		// TODO: Abort the read.
		return nil, err
	}

	length := binary.LittleEndian.Uint16(bhwTag[0:2])
	lengthCheck := binary.LittleEndian.Uint16(bhwTag[2:4])

	if (length|lengthCheck == 0) || (length^lengthCheck != 0xFFFF) {
		// Drop this packet...
		return nil, nil
	}

	if hwTag[0] == 12 /* && WLAN_IS_UP */ {
		var creditBuf [8]byte
		err = c.read(ioTransfer(creditBuf[:]))
		if err != nil {
			// TODO: Abort the read.
			return nil, err
		}
		c.updateCredit(creditBuf[:])
	}

	// Allocate a buffer to store the entire packet.
	data := make([]byte, length)

	// Copy data that was already read.
	n := copy(data, bhwTag)

	// Read the rest of the data.
	if int(length) > n {
		err = c.read(ioTransfer(data[n:]))
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (c *Cyw4343w[SDIO]) receiveOnePacket() error {
	data, err := c.readFrame()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		// The packet was likely dropped.
		return nil
	}

	// Process the packet.
	err = c.processRxPacket(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyw4343w[SDIO]) packetAvailableToRead() (bool, error) {
	base := sdiodCoreBaseAddress(c.chipId)

	// TODO: Check that the WLAN backplane is up before continuing...

	// Read the interrupt status.
	irqStatus, err := c.readBackplaneValue(base+0x20, 4)
	if err != nil {
		return false, err
	}

	if irqStatus&iHmbHostInt != 0 {
		// Read mailbox data and ack that we did so.
		hmbData, err := c.readBackplaneValue(base+0x4C, 4)
		if err == nil && hmbData > 0 {
			// Acknowledge.
			err = c.writeBackplaneValue(base+0x40, 4, smbIntAck)
			if err != nil {
				return false, err
			}
		}

		if hmbData&IHmbDataFwHalt != 0 {
			return false, hal.ErrInvalidState
		}

		if irqStatus&iHmbSwMask != 0 {
			// Clear any interrupts.
			err = c.writeBackplaneValue(base+0x20, 4, irqStatus&iHmbSwMask)
			if err != nil {
				return false, err
			}
		}
	}
	return irqStatus&iHmbSwMask != 0, nil
}
