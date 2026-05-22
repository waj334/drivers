package cyw4343w

import (
	"sync"

	"pkg.si-go.dev/chip/core/hal"
)

func (c *Cyw4343w[HostT, CacheT]) SetIovar(name string, data []byte) (BufferHandle, error) {
	// Allocate the buffer.
	buffer, payload, err := c.iovarBuffer(name, len(data))
	if err != nil {
		return BufferHandle{}, err
	}

	// Copy the input data into the payload.
	copy(payload.Data, data)

	// Write the iovar.
	response, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcSetVar,
		cmdType: cdcSet,
		data:    buffer.Data,
	})

	// Release the buffer back to the pool.
	buffer.Close()

	if err != nil {
		return BufferHandle{}, err
	}

	return response, nil
}

func (c *Cyw4343w[HostT, CacheT]) Iovar(name string, dataLen int) (BufferHandle, error) {
	// Allocate the buffer.
	buffer, _, err := c.iovarBuffer(name, dataLen)
	if err != nil {
		return BufferHandle{}, err
	}

	// Read the iovar from the device.
	response, err := c.sendIoctl(ioctlCommandType{
		data:    buffer.Data,
		cmdType: cdcGet,
		cmd:     wlcGetVar,
	})

	// Release the buffer back to the pool.
	buffer.Close()

	if err != nil {
		return BufferHandle{}, err
	}

	return response, nil
}

func (c *Cyw4343w[HostT, CacheT]) iovarBuffer(name string, dataLen int) (BufferHandle, BufferHandle, error) {
	nameLen := len(name) + 1
	totalLen := uintptr(roundUp(int(ioctlCommandHeaderLength)+nameLen+dataLen, 4))

	if totalLen > c.controlPool.SlotSize() {
		return BufferHandle{}, BufferHandle{}, hal.ErrInvalidBuffer
	}

	slot := c.controlPool.Get()
	if slot == nil {
		return BufferHandle{}, BufferHandle{}, hal.ErrInvalidBuffer
	}
	buf := slot[:totalLen]
	// Zero is important — leftover bytes from a previous use would corrupt
	// the iovar string (no null terminator, garbage payload).
	for i := range buf {
		buf[i] = 0
	}
	copy(buf[ioctlCommandHeaderLength:], name)

	data := buf[int(ioctlCommandHeaderLength)+nameLen:]

	// Use sync once to prevent either handle from returning the same buffer twice.
	var once sync.Once
	release := func() {
		once.Do(func() {
			c.controlPool.Put(slot)
		})
	}

	return BufferHandle{
			Data:    buf,
			release: release,
		}, BufferHandle{
			Data:    data,
			release: release,
		}, nil
}
