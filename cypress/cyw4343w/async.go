package cyw4343w

import (
	"unsafe"
)

const (
	etherTypeBrcm = 0x886C
	brcmOui       = "\x00\x10\x18"
)

func (c *Cyw4343w[HostT, CacheT]) processAsync(handle BufferHandle) error {
	data := handle.Data

	if len(data) < int(unsafe.Sizeof(bcdHeaderType{})) {
		handle.Close()
		return nil
	}

	bcd := (*bcdHeaderType)(unsafe.Pointer(&data[0]))
	eventOffset := int(unsafe.Sizeof(bcdHeaderType{})) + int(bcd.dataOffset)*4

	if eventOffset+int(unsafe.Sizeof(eventType{})) > len(data) {
		handle.Close()
		return nil
	}

	event := (*eventType)(unsafe.Pointer(&data[eventOffset]))

	event.eth.ethernetType = ntoh16(event.eth.ethernetType)
	if event.eth.ethernetType != etherTypeBrcm {
		handle.Close()
		return nil
	}

	if memcmp(unsafe.Pointer(unsafe.StringData(brcmOui)), unsafe.Pointer(&event.header.oui[0]), 3) != 0 {
		handle.Close()
		return nil
	}

	event.event.eventType = ntoh32(event.event.eventType)
	event.event.status = ntoh32(event.event.status)
	event.event.reason = ntoh32(event.event.reason)
	event.event.authType = ntoh32(event.event.authType)
	event.event.dataLength = ntoh32(event.event.dataLength)

	payloadOffset := eventOffset + int(unsafe.Sizeof(eventType{}))
	if payloadOffset+int(event.event.dataLength) > len(data) {
		handle.Close()
		return nil
	}

	whdEvent := &event.event

	if whdEvent.eventType == wlcEPskSup {
		whdEvent.status += wlcSupStatusOffset
		whdEvent.reason += wlcESupReasonOffset
	} else if whdEvent.eventType == wlcEPrune {
		whdEvent.reason += wlcESupReasonOffset
	} else if whdEvent.eventType == wlcEDisassoc || whdEvent.eventType == wlcEDeauth {
		whdEvent.status += wlcSupStatusOffset
		whdEvent.reason += wlcESupReasonOffset
	}

	ev := AsyncEvent{
		Type:          whdEvent.eventType,
		Status:        whdEvent.status,
		Reason:        whdEvent.reason,
		Auth:          whdEvent.authType,
		Data:          data,
		EventOffset:   eventOffset,
		PayloadOffset: payloadOffset,
		Handle:        handle,
	}

	if !c.events.Dispatch(ev) {
		handle.Close()
	}

	return nil
}
