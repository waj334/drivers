package cyw4343w

import (
	"sync"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

type DataTransfer struct {
	Data     []byte
	Address  uint32
	Function uint32
}

var (
	transferMutex sync.Mutex
)

func (c *Cyw4343w[HostT, CacheT]) write(transfer DataTransfer) error {
	return c.transfer(transfer, true)
}

func (c *Cyw4343w[HostT, CacheT]) read(transfer DataTransfer) error {
	return c.transfer(transfer, false)
}

func (c *Cyw4343w[HostT, CacheT]) transfer(transfer DataTransfer, write bool) error {
	var resp sdio.Response
	var err error

	direction := sdio.Read
	if write {
		direction = sdio.Write
	}

	// Only one goroutine should be using the SDIO interface at a time.
	transferMutex.Lock()

	// Send the packet.
	if len(transfer.Data) == 1 {
		args := sdio.CMD52Args{
			Address:   transfer.Address,
			Raw:       sdio.Normal,
			Function:  transfer.Function,
			ReadWrite: direction,
		}

		if write {
			args.Data = uint32(transfer.Data[0])
		}

		// Write a single byte using CMD52.
		resp, err = c.host.SendCommand(sdio.Command{
			Data:     transfer.Data,
			Class:    sdio.CMD52,
			Argument: args.Value(),
		})

		if err != nil {
			transferMutex.Unlock()
			return err
		}

		// Check the response for an error.
		r5 := sdio.R5(resp[0])
		if r5.Flags()&sdio.R5ErrorBits != 0 {
			transferMutex.Unlock()
			return sdio.ErrCommandFail
		}

		if !write {
			// Copy the data into the buffer.
			transfer.Data[0] = r5.Data()
		}

		transferMutex.Unlock()
		return nil
	} else if len(transfer.Data) >= 64 {
		// Write using CMD53 in block mode.
		err := c.transferBlocks(transfer, write)
		transferMutex.Unlock()
		return err
	}

	// Transfer 2-63 bytes using CMD53 in byte mode.
	resp, err = c.host.SendCommand(sdio.Command{
		Data:  transfer.Data,
		Class: sdio.CMD53,
		Argument: sdio.CMD53Args{
			Count:     uint32(len(transfer.Data)),
			Address:   transfer.Address,
			OpCode:    sdio.Incrementing,
			BlockMode: sdio.Bytes,
			Function:  transfer.Function,
			ReadWrite: direction,
		}.Value(),
	})

	if err != nil {
		transferMutex.Unlock()
		return err
	}

	// Check the response for an error.
	r5 := sdio.R5(resp[0])
	if r5.Flags()&sdio.R5ErrorBits != 0 {
		transferMutex.Unlock()
		return sdio.ErrCommandFail
	}

	transferMutex.Unlock()
	return nil
}

func (c *Cyw4343w[HostT, CacheT]) transferBlocks(transfer DataTransfer, write bool) error {
	var resp sdio.Response
	var err error
	cmd := sdio.Command{
		Class:     sdio.CMD53,
		BlockSize: blockSize,
	}

	args := sdio.CMD53Args{
		Count:     uint32(max(len(transfer.Data)/blockSize, 1)),
		Address:   transfer.Address,
		OpCode:    sdio.Incrementing,
		BlockMode: sdio.Blocks,
		Function:  transfer.Function,
		ReadWrite: sdio.Read,
	}

	if write {
		args.ReadWrite = sdio.Write
	}

	remaining := len(transfer.Data) % blockSize
	fullBlockBytes := len(transfer.Data) - remaining
	if fullBlockBytes > 0 {
		cmd.Data = transfer.Data[:fullBlockBytes]
		cmd.Argument = args.Value()
		resp, err = c.host.SendCommand(cmd)
		if err != nil {
			return err
		}

		// Check the response for an error.
		r5 := sdio.R5(resp[0])
		if r5.Flags()&sdio.R5ErrorBits != 0 {
			return sdio.ErrCommandFail
		}
	}

	if remaining > 0 {
		// Allocate a block to store the remaining bytes.
		b := make([]byte, blockSize)
		copy(b, transfer.Data[fullBlockBytes:])

		// Set up the command args for the final transfer.
		cmd.Data = b
		args.Count = 1
		if args.OpCode == sdio.Incrementing {
			args.Address += uint32(fullBlockBytes)
		}

		// Send the remaining block.
		cmd.Argument = args.Value()
		resp, err = c.host.SendCommand(cmd)
		if err != nil {
			return err
		}

		// Check the response for an error.
		r5 := sdio.R5(resp[0])
		if r5.Flags()&sdio.R5ErrorBits != 0 {
			return sdio.ErrCommandFail
		}

		if !write {
			copy(transfer.Data[fullBlockBytes:], b[:remaining])
		}
	}

	return nil
}
