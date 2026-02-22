package cyw4343w

import "unsafe"

const (
	etherTypeBrcm = 0x886C
	brcmOui       = "\x00\x10\x18"
)

func (c *Cyw4343w[SDIO]) processAsync(data []byte) error {
	header := (*bcdHeaderType)(unsafe.Pointer(&data[0]))
	event := (*eventType)(unsafe.Pointer(&data[header.dataOffset+1]))

	if event.eth.ethernetType != etherTypeBrcm {
		// This is not an event.
		return nil
	}

	if memcmp(unsafe.Pointer(unsafe.StringData(brcmOui)), unsafe.Pointer(&event.header.oui[0]), 3) != 0 {
		// This is not a broadcom event type.
		return nil
	}

	// Verify that the data length is correct.
	if event.event.dataLength > uint32(len(data)-(int(uintptr(unsafe.Pointer(event))-uintptr(unsafe.Pointer(header))))) {
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

	// TODO: Decide how events should be handled...

	return nil
}
