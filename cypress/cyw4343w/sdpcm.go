package cyw4343w

import (
	"fmt"
	"os"
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

func (c *Cyw4343w[HostT, CacheT]) updateCredit(data []byte) {
	var txSeqMax uint8
	header := (*sdpcmSwHeaderType)(unsafe.Pointer(&data[4]))
	if (header.channelAndFlags & 0x0F) < 3 {
		txSeqMax = header.busDataCredit
		oldMax := c.txMax
		if txSeqMax-c.txSeq > 0x40 {
			txSeqMax = c.txSeq + 2
		}
		c.txMax = txSeqMax
		if txSeqMax != oldMax && c.debug {
			fmt.Fprintf(os.Stdout, "[CREDIT] update: seq=%d oldMax=%d newMax=%d\n", c.txSeq, oldMax, txSeqMax)
		}
	}
}

// hasCredit returns true if the driver has TX credit available to send a packet.
func (c *Cyw4343w[HostT, CacheT]) hasCredit() bool {
	return c.txSeq != c.txMax
}

// waitForCredits waits until a TX credit is available. The background polling goroutine updates credits.
func (c *Cyw4343w[HostT, CacheT]) waitForCredits() error {
	if c.hasCredit() {
		return nil
	}

	if c.debug {
		fmt.Fprintf(os.Stdout, "[CREDIT] waiting: seq=%d max=%d\n", c.txSeq, c.txMax)
	}

	deadline := time.Now().Add(defaultTimeout)
	for {
		if c.hasCredit() {
			if c.debug {
				fmt.Fprintf(os.Stdout, "[CREDIT] acquired: seq=%d max=%d\n", c.txSeq, c.txMax)
			}
			return nil
		}

		if time.Now().After(deadline) {
			if c.debug {
				fmt.Fprintf(os.Stdout, "[CREDIT] TIMEOUT: seq=%d max=%d\n", c.txSeq, c.txMax)
			}
			return sdio.ErrTimeout
		}

		time.Sleep(time.Millisecond)
	}
}

func (c *Cyw4343w[HostT, CacheT]) processRxPacket(handle BufferHandle) error {
	header := *((*sdpcmHeaderType)(unsafe.Pointer(&handle.Data[0])))

	// Extract the total SDPCM packet size from the first two frametag bytes.
	size := header.frameTag[0]

	// Check that the second two frametag bytes are the binary inverse of the size.
	sizeInv := ^size
	if header.frameTag[1] != sizeInv {
		handle.Close()
		return hal.ErrInvalidBuffer
	}

	// Check whether the packet is big enough to contain the SDPCM header OR if it is too big to handle.
	if size < uint16(sdpcmHeaderLength) || int(size) > len(handle.Data) {
		handle.Close()
		return hal.ErrInvalidBuffer
	}

	if size == uint16(sdpcmHeaderLength) {
		// This is a flow control packet with no data.
		handle.Close()
		return nil
	}

	// Update credits.
	c.updateCredit(handle.Data)

	// Use the header_length field as the data offset — the firmware may include
	// extension headers or padding beyond the fixed 12-byte SDPCM header.
	dataOffset := uint16(header.headerLength)
	if dataOffset < uint16(sdpcmHeaderLength) {
		dataOffset = uint16(sdpcmHeaderLength)
	}
	if dataOffset >= size {
		// The header_length exceeds the packet — nothing to process.
		handle.Close()
		return nil
	}

	ch := header.channelAndFlags & 0x0F
	packet := BufferHandle{
		Data:    handle.Data[dataOffset:],
		release: handle.release,
	}

	switch ch {
	case controlHeader:
		return c.processIoctl(packet)
	case dataHeader:
		if c.debug {
			fmt.Fprintf(os.Stdout, "[POLL] rx data frame, len=%d\n", size)
		}
		return c.processData(packet)
	case asynceventHeader:
		if c.debug {
			fmt.Fprintf(os.Stdout, "[POLL] rx async event, len=%d\n", size)
		}
		return c.processAsync(packet)
	default:
		if c.debug {
			fmt.Fprintf(os.Stdout, "[POLL] rx unknown channel=%d, len=%d\n", ch, size)
		}
		packet.Close()
		return nil
	}
}
