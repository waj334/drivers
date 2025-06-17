package cyw4343w

func (c *Cyw4343w[SDIO]) SetIovar(name string, data []byte) ([]byte, error) {
	// Allocate the buffer.
	buffer, data := iovarBuffer(name, len(data))

	// Copy the input data into the payload.
	copy(buffer, data)

	// Write the iovar.
	response, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcSetVar,
		cmdType: cdcSet,
		data:    buffer,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Cyw4343w[SDIO]) Iovar(name string, dataLen int) ([]byte, error) {
	// Allocate the buffer.
	buffer, _ := iovarBuffer(name, dataLen)

	// Read the iovar from the device.
	response, err := c.sendIoctl(ioctlCommandType{
		data:    buffer,
		cmdType: cdcGet,
		cmd:     wlcGetVar,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func iovarBuffer(name string, dataLen int) (buffer []byte, data []byte) {
	nameLen := len(name) + 1 // +1 for null terminator
	totalLen := roundUp(int(ioctlCommandHeaderLength)+nameLen+dataLen, 4)

	buffer = make([]byte, totalLen)

	// Copy name in the buffer.
	copy(buffer[ioctlCommandHeaderLength:], name)

	// Data is immediately after the string.
	data = buffer[int(ioctlCommandHeaderLength)+nameLen:]

	return buffer, data
}
