package cyw4343w

import (
	"time"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

const (
	ioctlHeaderLength = unsafe.Sizeof(ioctlCommandHeaderType{})
)

type ioctlHeaderType struct {
	cmd    uint32
	length [2]uint16
	flags  uint32
	status uint32
}

type ioctlCommandHeaderType struct {
	sdpcmHeaderType
	ioctlHeaderType
}

type ioctlCommandType struct {
	cmd     uint32
	cmdType sdpcmCommandType
	data    []byte
}

var ioctlTransfer = sdio.Transfer{
	Address:   0,
	BlockSize: 64,
	Function:  uint8(wlanFunction),
	Increment: false,
}

func (c *Cyw4343w[SDIO]) sendIoctl(cmd ioctlCommandType) ([]byte, error) {
	c.iovarMutex.Lock()
	ioctlId := c.requestId()
	dataLength := uint16(len(cmd.data))
	totalLength := uint16(ioctlHeaderLength) + dataLength

	// Prepare the SDPCM Header.
	header := ioctlCommandHeaderType{
		sdpcmHeaderType: sdpcmHeaderType{
			frameTag: [2]uint16{
				totalLength,
				^totalLength,
			},
			sdpcmSwHeaderType: sdpcmSwHeaderType{
				sequence:        c.txSeq,
				channelAndFlags: controlHeader,
				headerLength:    uint8(sdpcmHeaderLength),
			},
		},
		ioctlHeaderType: ioctlHeaderType{
			cmd: cmd.cmd,
			length: [2]uint16{
				dataLength,
				0,
			},
			flags:  ((ioctlId << cdcfIocIdShift) & cdcfIocIdMask) | uint32(cmd.cmdType) | (c.bsscIndex << cdcfIocIfShift),
			status: 0,
		},
	}
	c.txSeq++

	// Allocate a buffer for the payload.
	headerData := toSlice(&header)
	data := make([]byte, len(headerData)+len(cmd.data))

	// Copy the header data followed by the command data to the buffer.
	n := copy(data, headerData)
	copy(data[n:], cmd.data)

	// Send the IOCTL command.
	err := c.write(data, ioctlTransfer)
	if err != nil {
		return nil, err
	}

	// Wait for the response.
	deadline := time.Now().Add(defaultTimeout)
	for {
		response, ok := c.receiveQueue.Dequeue(ioctlId)
		if ok {
			c.iovarMutex.Unlock()
			return response, nil
		} else if time.Now().After(deadline) {
			c.iovarMutex.Unlock()
			return nil, sdio.ErrTimeout
		}
	}
}

func (c *Cyw4343w[SDIO]) processIoctl(data []byte) error {
	header := *((*ioctlHeaderType)(unsafe.Pointer(&data[0])))
	id := (header.flags & cdcfIocIdMask) >> cdcfIocIdShift
	c.receiveQueue.Insert(id, data[ioctlHeaderLength:])
	return nil
}
