package cyw4343w

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
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

// abortRead sends an I/O abort to terminate a failed SDIO transfer on F2.
func (c *Cyw4343w[HostT, CacheT]) abortRead() {
	// Write to the CCCR I/O Abort register to abort the current F2 operation.
	_ = c.writeRegisterValue(busFunction, sdiodCccrIoabort, 1, uint32(wlanFunction))

	// Also set the read frame terminate bit in the frame control register.
	_ = c.writeRegisterValue(backplaneFunction, sdioFrameControl, 1, sfcRfTerm)
}

func (c *Cyw4343w[HostT, CacheT]) Poll() error {
	c.mutex.Lock()
	err := c.poll()
	c.mutex.Unlock()
	return err
}

func (c *Cyw4343w[HostT, CacheT]) poll() error {
	// Ensure the bus is awake before polling.
	if err := c.busWake(); err != nil {
		return err
	}

	// Check if the interrupt indicated there is a packet to read.
	available, err := c.packetAvailableToRead()
	if err != nil {
		return err
	}

	if available {
		err = c.receiveOnePacket()
		if err != nil {
			if c.debug {
				fmt.Fprintf(os.Stdout, "[POLL] receiveOnePacket error: %s\n", err.Error())
			}
			return err
		}
	}

	return nil
}

func (c *Cyw4343w[HostT, CacheT]) readFrame() (BufferHandle, error) {
	// Check that the WLAN backplane is up before continuing.
	up, err := c.ensureBackplaneUp()
	if err != nil {
		return BufferHandle{}, err
	}
	if !up {
		return BufferHandle{}, nil
	}

	var hwTag [2]uint16
	bhwTag := unsafe.Slice((*byte)(unsafe.Pointer(&hwTag[0])), 4)

	// Read the frame header and verify validity.
	err = c.read(ioTransfer(bhwTag[:4]))
	if err != nil {
		c.abortRead()
		return BufferHandle{}, err
	}

	length := binary.LittleEndian.Uint16(bhwTag[0:2])
	lengthCheck := binary.LittleEndian.Uint16(bhwTag[2:4])

	if (length|lengthCheck == 0) || (length^lengthCheck != 0xFFFF) {
		// Drop this packet...
		return BufferHandle{}, nil
	}

	if hwTag[0] == 12 && c.busIsUp {
		if c.debug {
			fmt.Fprintf(os.Stdout, "[POLL] credit-only packet\n")
		}
		var creditBuf [8]byte
		err = c.read(ioTransfer(creditBuf[:]))
		if err != nil {
			c.abortRead()
			return BufferHandle{}, err
		}
		c.updateCredit(creditBuf[:])
		//return nil, nil
	}

	if uintptr(length) > c.rxPool.SlotSize() {
		return BufferHandle{}, hal.ErrInvalidBuffer
	}

	data := c.rxPool.Get()
	if data == nil {
		return BufferHandle{}, errPoolExhausted
	}

	// Copy data that was already read.
	n := copy(data, bhwTag)

	// Read the rest of the data.
	if int(length) > n {
		c.rxPool.CompleteRx(data[n:length])
		err = c.read(ioTransfer(data[n:length]))
		if err != nil {
			c.rxPool.Put(data)
			return BufferHandle{}, err
		}

		c.rxPool.CompleteRx(data[n:length])
	}

	var once sync.Once
	return BufferHandle{
		Data: data,
		release: func() {
			once.Do(func() { c.rxPool.Put(data) })
		},
	}, nil
}

func (c *Cyw4343w[HostT, CacheT]) receiveOnePacket() error {
	handle, err := c.readFrame()
	if err != nil {
		return err
	}

	if len(handle.Data) == 0 {
		// The packet was likely dropped.
		return nil
	}

	// Process the packet.
	err = c.processRxPacket(handle)
	return err
}

func (c *Cyw4343w[HostT, CacheT]) packetAvailableToRead() (bool, error) {
	base := sdiodCoreBaseAddress(c.chipId)

	// Check that the WLAN backplane is up before continuing.
	up, err := c.ensureBackplaneUp()
	if err != nil {
		return false, err
	}
	if !up {
		return false, nil
	}

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
