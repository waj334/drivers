package cyw4343w

import (
	"fmt"
	"net"
	"os"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal"
)

// SetRxCallback registers a function that will be called when an Ethernet frame
// is received from the WiFi chip on the data channel. The frame passed to fn is
// a raw Ethernet frame (dest MAC | src MAC | ethertype | payload).
func (c *Cyw4343w[HostT, CacheT]) SetRxCallback(fn func(frame net.RxFrame)) {
	c.rxCallback = fn
}

// processData handles incoming data frames from SDPCM channel 2.
func (c *Cyw4343w[HostT, CacheT]) processData(handle BufferHandle) error {
	if c.rxCallback == nil {
		handle.Close()
		return nil
	}

	if len(handle.Data) < int(unsafe.Sizeof(bcdHeaderType{})) {
		handle.Close()
		return nil
	}

	// Parse the BCD header to find where the Ethernet frame begins.
	header := (*bcdHeaderType)(unsafe.Pointer(&handle.Data[0]))
	offset := int(unsafe.Sizeof(bcdHeaderType{})) + int(header.dataOffset)*4

	if offset >= len(handle.Data) {
		handle.Close()
		return nil
	}

	frame := net.RxFrame{
		Data:    handle.Data[offset:],
		Release: handle.release,
	}

	if c.debug {
		if len(frame.Data) >= 14 {
			fmt.Fprintf(os.Stdout, "[RX] len=%d dst=%02x:%02x:%02x:%02x:%02x:%02x src=%02x:%02x:%02x:%02x:%02x:%02x type=%02x%02x\n",
				len(frame.Data),
				frame.Data[0], frame.Data[1], frame.Data[2], frame.Data[3], frame.Data[4], frame.Data[5],
				frame.Data[6], frame.Data[7], frame.Data[8], frame.Data[9], frame.Data[10], frame.Data[11],
				frame.Data[12], frame.Data[13])
		}
	}

	c.rxCallback(frame)
	return nil
}

// SendEthernet transmits a raw Ethernet frame over the WiFi data channel.
func (c *Cyw4343w[HostT, CacheT]) SendEthernet(frame []byte) error {
	const bcdLen = 4 // sizeof(bcdHeaderType)
	totalLen := uint16(sdpcmHeaderLength) + bcdLen + uint16(len(frame))
	if uintptr(totalLen) > c.txPool.SlotSize() {
		return hal.ErrInvalidBuffer
	}

	if c.debug {
		if len(frame) >= 14 {
			fmt.Fprintf(os.Stdout, "[TX] len=%d dst=%02x:%02x:%02x:%02x:%02x:%02x src=%02x:%02x:%02x:%02x:%02x:%02x type=%02x%02x seq=%d max=%d\n",
				len(frame),
				frame[0], frame[1], frame[2], frame[3], frame[4], frame[5],
				frame[6], frame[7], frame[8], frame[9], frame[10], frame[11],
				frame[12], frame[13],
				c.txSeq, c.txMax)
		}
	}

	c.iovarMutex.Lock()

	if err := c.busWake(); err != nil {
		if c.debug {
			fmt.Fprintf(os.Stdout, "[TX] busWake error: %s\n", err.Error())
		}
		c.iovarMutex.Unlock()
		return err
	}

	if err := c.waitForCredits(); err != nil {
		if c.debug {
			fmt.Fprintf(os.Stdout, "[TX] waitForCredits timeout: seq=%d max=%d\n", c.txSeq, c.txMax)
		}
		c.iovarMutex.Unlock()
		return err
	}

	buf := c.txPool.Get()
	if buf == nil {
		return errPoolExhausted
	}

	pkt := buf[:totalLen]

	// Build the SDPCM header.
	header := (*sdpcmHeaderType)(unsafe.Pointer(&buf[0]))
	*header = sdpcmHeaderType{
		frameTag: [2]uint16{
			totalLen,
			^totalLen,
		},
		sdpcmSwHeaderType: sdpcmSwHeaderType{
			sequence:        c.txSeq,
			channelAndFlags: dataHeader,
			headerLength:    uint8(sdpcmHeaderLength),
		},
	}

	// BDC header: set protocol version 2 in flags byte; priority, flags2,
	// and dataOffset remain zero.
	pkt[sdpcmHeaderLength] = bdcFlagVersion
	copy(pkt[sdpcmHeaderLength+uintptr(bcdLen):], frame)
	c.txSeq++

	c.txPool.PrepareTx(pkt)
	err := c.write(ioTransfer(pkt))

	// Release the buffer back to the pool
	c.txPool.Put(buf)

	if c.debug {
		if err != nil {
			fmt.Fprintf(os.Stdout, "[TX] write error: %s\n", err.Error())
		} else {
			fmt.Fprintf(os.Stdout, "[TX] sent ok, newSeq=%d\n", c.txSeq)
		}
	}
	c.iovarMutex.Unlock()
	return err
}
