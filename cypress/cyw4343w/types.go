package cyw4343w

const msgIfaceNameMax = 16

type macType [6]byte

type ethernetHeaderType struct {
	destination  macType
	source       macType
	ethernetType uint16
}

type ethernetEventHeaderType struct {
	subtype     uint16
	length      uint16
	version     uint8
	oui         [3]uint8
	userSubtype uint16
}

type eventMessageType struct {
	version        uint16
	flags          uint16
	eventType      uint32
	status         uint32
	reason         uint32
	authType       uint32
	dataLength     uint32
	addr           macType
	ifaceName      [msgIfaceNameMax]byte
	ifaceIndex     uint8
	bssConfigIndex uint8
}

type eventType struct {
	eth    ethernetHeaderType
	header ethernetEventHeaderType
	event  eventMessageType
	/* Data is below these fields */
}

type bcdHeaderType struct {
	flags      byte
	priority   byte
	flags2     byte
	dataOffset byte
}

type BufferHandle struct {
	Data    []byte
	release func()
}

func (b BufferHandle) Close() {
	b.release()
}
