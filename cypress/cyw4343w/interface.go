package cyw4343w

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

const (
	ioctlTypeGet = 0x00
	ioctlTypeSet = 0x02

	busFunction       = 0
	backplaneFunction = 1
	wlanFunction      = 2

	defaultTimeout = 100 * time.Millisecond
)

type Config[Host sdio.Host] struct {
	Host Host
}

type Cyw4343w[Host sdio.Host] struct {
	host         Host
	receiveQueue queue

	requestIdCounter uintptr
	bsscIndex        int // Leaving at 0 for now.

	mutex      sync.Mutex
	iovarMutex sync.Mutex
}

func New[Host sdio.Host]() *Cyw4343w[Host] {
	return &Cyw4343w[Host]{}
}

func (c *Cyw4343w[SDIO]) Configure(config Config[SDIO]) error {
	c.mutex.Lock()
	c.host = config.Host
	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) Initialize() error {
	c.mutex.Lock()
	var resp sdio.Response
	var err error
	ready := false

	// Send CMD0 to reset the card
	if _, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD0, Argument: 0}); err != nil {
		c.mutex.Unlock()
		return err
	}

	for retry := 0; retry < 1000; retry++ {
		// Send CMD5 to get the ready status of the card.
		if resp, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD5, Argument: 0x00FF8000}); err != nil {
			time.Sleep(time.Millisecond * 5)
			continue
		} else {
			if resp[0]>>31 == 0 {
				// The card is not ready. Try again.
				time.Sleep(time.Millisecond * 5)
				continue
			}
			ready = true
		}
		break
	}

	if !ready {
		c.mutex.Unlock()
		return sdio.ErrNotReady
	}

	// Send CMD3 to get the address of the card.
	resp, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD3, Argument: 0})
	if err != nil {
		c.mutex.Unlock()
		return err
	}

	// Send CMD7 with the returned RCA to select the card.
	_, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD7, Argument: resp[0] & 0xFFFF0000})
	if err != nil {
		c.mutex.Unlock()
		return err
	}

	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) Poll() error {
	c.mutex.Lock()
	// TODO: Receive packets from the device and place them in the receive queue.
	c.mutex.Unlock()
	return nil
}

func (c *Cyw4343w[SDIO]) requestId() uintptr {
	return atomic.AddUintptr(&c.requestIdCounter, 1)
}

type ioctlCommand struct {
	buf     []byte
	cmd     uint16
	cmdType uint8
}

func (c *Cyw4343w[SDIO]) sendIoctl(cmd ioctlCommand) ([]byte, error) {
	// Verify the command.
	if cmd.cmd > wlcLast {
		return nil, errInvalidCommand
	}

	c.iovarMutex.Lock()
	requestId := c.requestId()
	dataLen := len(cmd.buf) - int(unsafe.Sizeof(commonBusHeader{})) - int(unsafe.Sizeof(cdcHeader{}))
	sendPacket := (*controlHeader)(unsafe.Pointer(unsafe.SliceData(cmd.buf)))

	// TODO: IOVAR creates an unaligned data section. The original has a compensation for this.

	// Set up the packet this will be sent to the device.
	// TODO: Might need byte swapping to account for differences in endianness.
	sendPacket.cmd = uint32(cmd.cmd)
	sendPacket.len = uint32(dataLen)
	sendPacket.flags = uint32((requestId<<cdcfIocIdShift)&cdcfIocIdMask) |
		uint32(cmd.cmdType) | uint32(c.bsscIndex<<cdcfIocIfShift)

	// Send the packet.
	err := c.host.WriteBlocks(cmd.buf, sdio.Transfer{
		Address:   0,
		BlockSize: 64,
		Function:  wlanFunction,
		Increment: true,
	})

	if err != nil {
		c.iovarMutex.Unlock()
		return nil, err
	}

	// Wait for the response.
	deadline := time.Now().Add(defaultTimeout)
	for {
		response, ok := c.receiveQueue.Dequeue(requestId)
		if ok {
			c.iovarMutex.Unlock()
			return response, nil
		} else if time.Now().After(deadline) {
			c.iovarMutex.Unlock()
			return nil, sdio.ErrTimeout
		}
	}
}

func (c *Cyw4343w[SDIO]) Iovar(name string, dataLen int) ([]byte, error) {
	// Allocate the buffer.
	buffer, _ := iovarBuffer(name, dataLen)

	// Read the iovar from the device.
	response, err := c.sendIoctl(ioctlCommand{
		buf:     buffer,
		cmdType: ioctlTypeGet,
		cmd:     wlcGetVar,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func iovarBuffer(name string, dataLen int) (buffer []byte, data []byte) {
	// NOTE: Pad extra byte at the end of the name to account for null terminator that Go strings will not have.
	nameLen := len(name) + 1
	nameLenAlignmentOffset := (64 - nameLen) % int(unsafe.Sizeof(uintptr(0)))
	sz := ioctlOffset + nameLen + nameLenAlignmentOffset + dataLen

	// Allocate the buffer.
	buffer = make([]byte, sz)

	// Perform unsafe string to slice conversion to avoid heap allocation.
	ptr := unsafe.StringData(name)
	s := unsafe.Slice(ptr, len(name))

	// Copy the name into the beginning of the buffer.
	copy(buffer, s)

	// Create a slice at the start of where data can be read.
	data = buffer[nameLen+nameLenAlignmentOffset:]

	// Return the base buffer slice and a slice at the start of the data.
	return buffer, data
}
