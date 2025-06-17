package cyw4343w

import (
	"pkg.si-go.dev/chip/core/hal"
	"unsafe"
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
	// TODO: Choose to enable flow control here...
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

	// Check the SDPCM channel to decide what to do with the packet.
	switch header.channelAndFlags & 0x0F {
	case controlHeader:
		return c.processIoctl(data[sdpcmHeaderLength:])
	case dataHeader:
	case asynceventHeader:
	default:
		return hal.ErrInvalidState
	}

	c.updateCredit(data)

	return nil
}
