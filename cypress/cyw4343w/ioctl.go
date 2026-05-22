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

func ioctlPacket[T any](value T) ([]byte, *T) {
	buf := make([]byte, int(ioctlCommandHeaderLength)+int(unsafe.Sizeof(value)))
	return buf, (*T)(unsafe.Pointer(&buf[ioctlCommandHeaderLength]))
}

func (c *Cyw4343w[HostT, CacheT]) sendIoctl(cmd ioctlCommandType) (BufferHandle, error) {
	if len(cmd.data) < int(ioctlCommandHeaderLength) {
		return BufferHandle{}, errInvalidCommand
	}

	c.iovarMutex.Lock()

	if err := c.busWake(); err != nil {
		c.iovarMutex.Unlock()
		return BufferHandle{}, err
	}

	if err := c.waitForCredits(); err != nil {
		c.iovarMutex.Unlock()
		return BufferHandle{}, err
	}

	ioctlId := c.requestId()
	totalLength := uint16(len(cmd.data))
	dataLength := totalLength - uint16(ioctlCommandHeaderLength)

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

	c.controlPool.PrepareTx(cmd.data)
	err := c.write(ioctlTransfer(cmd.data))
	if err != nil {
		c.iovarMutex.Unlock()
		return BufferHandle{}, err
	}

	deadline := time.Now().Add(defaultTimeout * 2)
	for {
		responseHandle, ok := c.receiveQueue.Dequeue(ioctlId)
		if ok {
			response := responseHandle.Data
			header := *(*ioctlHeaderType)(unsafe.Pointer(&response[0]))

			c.iovarMutex.Unlock()

			if int32(header.status) != 0 {
				responseHandle.Close()
				return BufferHandle{}, errIoctlFailed
			}

			return BufferHandle{
				Data:    response[ioctlHeaderLength:],
				release: responseHandle.release,
			}, nil
		}

		if time.Now().After(deadline) {
			c.iovarMutex.Unlock()
			return BufferHandle{}, sdio.ErrTimeout
		}

		time.Sleep(time.Millisecond)
	}
}

func (c *Cyw4343w[HostT, CacheT]) processIoctl(handle BufferHandle) error {
	header := *((*ioctlHeaderType)(unsafe.Pointer(&handle.Data[0])))
	id := (header.flags & cdcfIocIdMask) >> cdcfIocIdShift
	c.receiveQueue.Insert(id, handle)
	return nil
}
