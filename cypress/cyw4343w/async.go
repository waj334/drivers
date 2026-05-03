package cyw4343w

import "unsafe"

const (
	etherTypeBrcm = 0x886C
	brcmOui       = "\x00\x10\x18"
)

func (c *Cyw4343w[SDIO]) processAsync(data []byte) error {
	header := (*bcdHeaderType)(unsafe.Pointer(&data[0]))
	eventOffset := int(unsafe.Sizeof(bcdHeaderType{})) + int(header.dataOffset)*4
	if eventOffset+int(unsafe.Sizeof(eventType{})) > len(data) {
		return nil
	}
	event := (*eventType)(unsafe.Pointer(&data[eventOffset]))

	// The ethernet and event header fields are in network byte order (big-endian).
	// Byte-swap them to host order in-place so all downstream readers see native values.
	event.eth.ethernetType = ntoh16(event.eth.ethernetType)

	if event.eth.ethernetType != etherTypeBrcm {
		// This is not an event.
		return nil
	}

	if memcmp(unsafe.Pointer(unsafe.StringData(brcmOui)), unsafe.Pointer(&event.header.oui[0]), 3) != 0 {
		// This is not a broadcom event type.
		return nil
	}

	// Swap event message fields from network to host order.
	event.event.eventType = ntoh32(event.event.eventType)
	event.event.status = ntoh32(event.event.status)
	event.event.reason = ntoh32(event.event.reason)
	event.event.authType = ntoh32(event.event.authType)
	event.event.dataLength = ntoh32(event.event.dataLength)

	// Verify that the data length is correct.
	if event.event.dataLength > uint32(len(data)-eventOffset-int(unsafe.Sizeof(eventType{}))) {
		// The size does not match.
		return nil
	}

	whdEvent := &event.event

	/* This is necessary because people who defined event statuses and reasons overlapped values. */
	if whdEvent.eventType == wlcEPskSup {
		whdEvent.status += wlcSupStatusOffset
		whdEvent.reason += wlcESupReasonOffset
	} else if whdEvent.eventType == wlcEPrune {
		whdEvent.reason += wlcESupReasonOffset
	} else if whdEvent.eventType == wlcEDisassoc || whdEvent.eventType == wlcEDeauth {
		whdEvent.status += wlcSupStatusOffset
		whdEvent.reason += wlcESupReasonOffset
	}

	// Queue the event for waiting callers.
	c.eventQueue.Insert(whdEvent.eventType, data)

	return nil
}
