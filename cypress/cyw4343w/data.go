package cyw4343w

import (
	"fmt"
	"os"
	"unsafe"
)

// SetRxCallback registers a function that will be called when an Ethernet frame
// is received from the WiFi chip on the data channel. The frame passed to fn is
// a raw Ethernet frame (dest MAC | src MAC | ethertype | payload).
func (c *Cyw4343w[SDIO]) SetRxCallback(fn func([]byte)) {
	c.rxCallback = fn
}

// processData handles incoming data frames from SDPCM channel 2.
func (c *Cyw4343w[SDIO]) processData(data []byte) error {
	if c.rxCallback == nil {
		return nil
	}

	if len(data) < int(unsafe.Sizeof(bcdHeaderType{})) {
		return nil
	}

	// Parse the BCD header to find where the Ethernet frame begins.
	header := (*bcdHeaderType)(unsafe.Pointer(&data[0]))
	offset := int(unsafe.Sizeof(bcdHeaderType{})) + int(header.dataOffset)*4

	if offset >= len(data) {
		return nil
	}

	frame := data[offset:]
	if c.debug {
		if len(frame) >= 14 {
			fmt.Fprintf(os.Stdout, "[RX] len=%d dst=%02x:%02x:%02x:%02x:%02x:%02x src=%02x:%02x:%02x:%02x:%02x:%02x type=%02x%02x\n",
				len(frame),
				frame[0], frame[1], frame[2], frame[3], frame[4], frame[5],
				frame[6], frame[7], frame[8], frame[9], frame[10], frame[11],
				frame[12], frame[13])
		}
	}

	c.rxCallback(frame)
	return nil
}

// SendEthernet transmits a raw Ethernet frame over the WiFi data channel.
func (c *Cyw4343w[SDIO]) SendEthernet(frame []byte) error {
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

	const bcdLen = 4 // sizeof(bcdHeaderType)
	totalLen := uint16(sdpcmHeaderLength) + bcdLen + uint16(len(frame))

	buf := make([]byte, totalLen)

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
	buf[sdpcmHeaderLength] = bdcFlagVersion

	// Copy the Ethernet frame after the BCD header.
	copy(buf[uint16(sdpcmHeaderLength)+bcdLen:], frame)

	c.txSeq++

	err := c.write(ioTransfer(buf))
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
