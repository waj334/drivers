package cyw4343w

import (
	"time"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal"
	"pkg.si-go.dev/chip/core/hal/sdio"
)

const (
	sdpcmHeaderLength                  = unsafe.Sizeof(sdpcmHeaderType{})
	cdcGet            sdpcmCommandType = 0x00
	cdcSet            sdpcmCommandType = 0x02
)

type sdpcmCommandType uint32

type sdpcmHeaderType struct {
	frameTag [2]uint16
	sdpcmSwHeaderType
}

type sdpcmSwHeaderType struct {
	sequence            uint8
	channelAndFlags     uint8
	nextLength          uint8
	headerLength        uint8
	wirelessFlowControl uint8
	busDataCredit       uint8
	_                   [2]uint8
}

func (c *Cyw4343w[SDIO]) updateCredit(data []byte) {
	var txSeqMax uint8
	header := (*sdpcmSwHeaderType)(unsafe.Pointer(&data[4]))
	if (header.channelAndFlags & 0x0F) < 3 {
		txSeqMax = header.busDataCredit
		if txSeqMax-c.txSeq > 0x40 {
			txSeqMax = c.txSeq + 2
		}
		c.txMax = txSeqMax
	}
}

// hasCredit returns true if the driver has TX credit available to send a packet.
func (c *Cyw4343w[SDIO]) hasCredit() bool {
	return c.txSeq != c.txMax
}

// waitForCredits waits until a TX credit is available. The background polling goroutine updates credits.
func (c *Cyw4343w[SDIO]) waitForCredits() error {
	if c.hasCredit() {
		return nil
	}

	deadline := time.Now().Add(defaultTimeout)
	for {
		if c.hasCredit() {
			return nil
		}

		if time.Now().After(deadline) {
			return sdio.ErrTimeout
		}

		time.Sleep(time.Millisecond)
	}
}

func (c *Cyw4343w[SDIO]) processRxPacket(data []byte) error {
	header := *((*sdpcmHeaderType)(unsafe.Pointer(&data[0])))

	// Extract the total SDPCM packet size from the first two frametag bytes.
	size := header.frameTag[0]

	// Check that the second two frametag bytes are the binary inverse of the size.
	sizeInv := ^size
	if header.frameTag[1] != sizeInv {
		return hal.ErrInvalidBuffer
	}

	// Check whether the packet is big enough to contain the SDPCM header OR if it is too big to handle.
	if size < uint16(sdpcmHeaderLength) || int(size) > len(data) {
		return hal.ErrInvalidBuffer
	}

	if size == uint16(sdpcmHeaderLength) {
		// This is a flow control packet with no data.
		return nil
	}

	// Update credits.
	c.updateCredit(data)

	// Use the header_length field as the data offset — the firmware may include
	// extension headers or padding beyond the fixed 12-byte SDPCM header.
	dataOffset := uint16(header.headerLength)
	if dataOffset < uint16(sdpcmHeaderLength) {
		dataOffset = uint16(sdpcmHeaderLength)
	}
	if dataOffset >= size {
		// The header_length exceeds the packet — nothing to process.
		return nil
	}
	packet := data[dataOffset:]
	switch header.channelAndFlags & 0x0F {
	case controlHeader:
		return c.processIoctl(packet)
	case dataHeader:
		return c.processData(packet)
	case asynceventHeader:
		return c.processAsync(packet)
	default:
		// Silently ignore unhandled channels (e.g. glom/aggregation).
		return nil
	}
}
