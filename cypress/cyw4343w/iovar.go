package cyw4343w

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

/*
func iovarBuffer(name string, dataLen int) (buffer []byte, data []byte) {
	nameLen := len(name) + 1 // +1 for null terminator
	nameAlign := (64 - nameLen) % int(unsafe.Sizeof(uint32(0)))

	// Total size = IOCTL prefix + padding + name + data
	totalLen := ioctlOffset + nameAlign + nameLen + dataLen

	buffer = make([]byte, totalLen)

	// Where the name starts
	nameStart := ioctlOffset + nameAlign

	// Copy the name bytes
	copy(buffer[nameStart:], name)

	// Append null terminator
	buffer[nameStart+len(name)] = 0

	// Finally, the data region is after name+null
	data = buffer[nameStart+nameLen:]

	return buffer, data
}
*/

func iovarBuffer(name string, dataLen int) (buffer []byte, data []byte) {
	nameLen := len(name) + 1 // +1 for null terminator
	totalLen := roundUp(nameLen+dataLen, 4)

	buffer = make([]byte, totalLen)

	// Copy name in the buffer.
	copy(buffer, name)

	// Data is immediately after the string.
	data = buffer[nameLen:]

	return buffer, data
}
