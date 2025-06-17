package cyw4343w

import (
	"time"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

const (
	ioctlCommandHeaderLength = unsafe.Sizeof(ioctlCommandHeaderType{})
	ioctlHeaderLength        = unsafe.Sizeof(ioctlHeaderType{})
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

func ioctlTransfer(data []byte) DataTransfer {
	return DataTransfer{
		Data:     data,
		Address:  0,
		Function: uint32(wlanFunction),
	}
}

func (c *Cyw4343w[SDIO]) sendIoctl(cmd ioctlCommandType) ([]byte, error) {
	if len(cmd.data) < int(ioctlCommandHeaderLength) {
		return nil, errInvalidCommand
	}

	c.iovarMutex.Lock()
	ioctlId := c.requestId()
	totalLength := uint16(len(cmd.data))

	// NOTE: The data parameter already includes memory to store the header. So, subtract the length of the header to
	//       derive the length of the data payload.
	dataLength := totalLength - uint16(ioctlCommandHeaderLength)

	// Prepare the SDPCM Header.
	header := (*ioctlCommandHeaderType)(unsafe.Pointer(&cmd.data[0]))
	*header = ioctlCommandHeaderType{
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

	// Send the IOCTL command.
	err := c.write(ioctlTransfer(cmd.data))
	if err != nil {
		c.iovarMutex.Unlock()
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
